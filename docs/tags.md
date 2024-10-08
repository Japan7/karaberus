# Tag system

This document serves as a description of the tags in karaberus.
(Ideally the tags should be described in the source code and the descriptions could
be served by karaberus itself, but it is too early for that)

# Types of tags

There are 5 types of tags:

* Author tags: tags related to the author of the karaoke.
* Artist tags: tags related to the artists of the original song.
* Media tags: tags related to the source media of the song.
* Video tags: tags related to the video of the karaoke.
* Audio tags: tags related to the audio of the karaoke.

Generic and media tags are user defined while video and music tags are
statically defined in karaberus. 

## Author tags

Tags related to the author of the karaoke timing.

## Artist tags

Tags related to the artists involved in the song's creation (not the video, but could be added if needed?)

Artists tags support alternative names/aliases.

## Media tags

Tags related to the source media.

Media tags have:
* A title
* A type of media (Anime, Game, Live action, Cartoon)
* Alternative titles if needed


## Video tags

Tags related to the video.

Video tags currently defined:
* **Music Video**: The video is specifically made for the song
* **Fanmade**: The video is fanmade
* **Stream**: The video comes from a stream
* **Concert**: The video is taken from a live concert
* **Advertisement**: The video is an ad
* **NSFW**: The video is not safe for work
* **SPOILER**: The video contains spoilers
* **EPILEPSY**: The video could cause a seizure to people with photosensitive epilepsy


## Audio tags

Tags related to the audio.

Audio tags currently defined:
* **Opening**: The audio is an opening song of the source media
* **Ending**: The audio is an ending song of the source media
* **Insert**: The audio is an insert song of the source media
* **Image song**: The audio is an image song of the source media
* **Live**: The audio is a live interpretation of the song
* **Remix/Cover**: The audio is a remix or cover of the original song


# Transition from dakaraneko

* Episode tag (EP) → Version: it is more obvious in the interface anyway
* VIDEO → Alternative title
* VTITLE → Alternative title
* VIDEO → Comment (do we need a source field?)
* OARTIST → Artist (is it something we want to keep?)
* AMV details tag: feed into media source
* AMV tag → Fanmade + Music Video
* LONG/COURT → Version: it's always relative to another version of the song
