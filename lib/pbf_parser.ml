let rec iter_channel ic =
  match input_char ic with
  | exception End_of_file -> ()
  | b0 -> iter_channel ic

let parse filename =
  let ic = open_in_bin filename in
  Fun.protect
    ~finally:(fun () -> close_in_noerr ic)
    (fun () -> iter_channel ic)
