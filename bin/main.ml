open Gpath

let () =
  match Pbf_parser.parse "gent.osm.pbf" with
  | Ok () -> ()
  | Error s -> Printf.eprintf "%s\n" s
