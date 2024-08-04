import { useNavigate } from "@solidjs/router";
import MugenImport from "../../components/MugenImport";

export default function KaraokeMugen() {
  const navigate = useNavigate();

  return (
    <>
      <h1 class="header">Karaoke Mugen Import</h1>

      <MugenImport
        onImport={(kara) => navigate("/karaoke/browse/" + kara.KaraID)}
      />
    </>
  );
}
