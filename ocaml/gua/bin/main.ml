open Lwt
open Cohttp
open Cohttp_lwt_unix

let usage_msg = "guthub-user-activity <username>"
let username = ref ""
let get_username un = username := un
let () = Arg.parse [] get_username usage_msg
let url = "https://api.github.com/users/" ^ !username ^ "/events"
let () = Printf.printf "url: %s\n" url

let body =
  Client.get (Uri.of_string url) >>= fun (resp, body) ->
  let code = resp |> Response.status |> Code.code_of_status in
  Printf.printf "Response code: %d\n" code;
  Printf.printf "Received a response from the server";
  Printf.printf "Headers: %s\n" (resp |> Response.headers |> Header.to_string);
  body |> Cohttp_lwt.Body.to_string >|= fun body ->
  Printf.printf "Body of length: %d\n" (String.length body);
  body

let () =
  let body = Lwt_main.run body in
  body |> String.length |> print_int;
  print_endline "Received body\n"

(* let () = *)
(*   let uri = Uri.of_string url in *)
