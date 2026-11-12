exception Parse_error of string

let ( let* ) = Result.bind
let fail fmt = Printf.ksprintf (fun s -> raise (Parse_error s)) fmt

let read_or_err ic n what =
  In_channel.really_input_string ic n
  |> Option.to_result ~none:(Printf.sprintf "Unexpected EOF on %s" what)

let int32_of_bytes b0 b1 b2 b3 =
  (Char.code b0 lsl 24)
  lor (Char.code b1 lsl 16)
  lor (Char.code b2 lsl 8)
  lor Char.code b3

let decode_blob_header bytes =
  String.iter (fun c -> Printf.printf "%d: %c\n" (int_of_char c) c) bytes

let read_header ic =
  match In_channel.input_char ic with
  | None -> Ok ()
  | Some b0 ->
      let* rest = read_or_err ic 3 "BlobHeader length" in
      let header_len = int32_of_bytes b0 rest.[0] rest.[1] rest.[2] in
      let* header_data = read_or_err ic header_len "BlobHeader" in
      let _indexdata = decode_blob_header header_data in
      Printf.printf "header_len: %d\n" header_len;
      Ok ()

let iter_channel ic =
  read_header ic

let parse filename =
  let ic = open_in_bin filename in
  Fun.protect ~finally:(fun () -> close_in_noerr ic) (fun () -> iter_channel ic)
