use async_std::task;
use interprocess::local_socket::{prelude::*, GenericFilePath, GenericNamespaced, Stream};
use serde::{Serialize, Serializer};
use std::collections::HashMap;
use std::io::{prelude::*, BufReader};
use std::time::Duration;

#[derive(Serialize)]
pub struct LoadFile {
    name: String,
    pub url: String,
    pub flags: String,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub index: Option<i32>,
    #[serde(skip_serializing_if = "HashMap::is_empty")]
    #[serde(serialize_with = "options_serializer")]
    pub options: HashMap<String, String>,
}

impl Default for LoadFile {
    fn default() -> LoadFile {
        LoadFile {
            name: "loadfile".to_string(),
            url: String::new(),
            flags: "replace".to_string(),
            index: None,
            options: HashMap::new(),
        }
    }
}

#[derive(Serialize)]
#[serde(untagged)]
enum MpvCommand {
    LoadFile(LoadFile),
}

#[derive(Serialize)]
struct MpvCommandWrapper {
    command: MpvCommand,
}

fn options_serializer<S>(map: &HashMap<String, String>, serializer: S) -> Result<S::Ok, S::Error>
where
    S: Serializer,
{
    let mut options_string = String::new();
    let mut first = true;
    for (key, value) in map.into_iter() {
        if !first {
            options_string.push_str(",");
        }
        options_string.push_str(format!("{key}=%{len}%{value}", len = value.len()).as_str());
        first = false;
    }
    return serializer.serialize_str(options_string.as_str());
}

#[derive(Default, Clone)]
pub struct Mpv {
    pub socket: String,
}

impl Mpv {
    pub async fn loadfile(&self, command: LoadFile) {
        self.run_command(MpvCommand::LoadFile(command)).await;
    }

    async fn run_command(&self, command: MpvCommand) {
        let socket = self.socket.clone();
        let name = if cfg!(windows) {
            socket.to_ns_name::<GenericNamespaced>()
        } else {
            socket.to_fs_name::<GenericFilePath>()
        }
        .unwrap();

        let command = MpvCommandWrapper { command };
        let mpv_command = format!("{}\n", serde_json::to_string(&command).unwrap());

        //TODO: Handle possible failures and response received
        let mut n_tries = 5;
        let conn = loop {
            match Stream::connect(name.clone()) {
                Ok(conn) => break conn,
                Err(err) => {
                    if n_tries > 0 {
                        n_tries -= 1;
                        task::sleep(Duration::from_millis(200)).await;
                    } else {
                        eprintln!("Unable to connect to mpv socket: {}", err);
                        return;
                    }
                }
            }
        };
        let mut conn = BufReader::new(conn);
        let mut buffer = String::with_capacity(1024);

        let _ = conn.get_mut().write_all(mpv_command.as_bytes());
        conn.read_line(&mut buffer).unwrap();
    }
}
