open Core
open Async
open Lwt
open Cohttp
open Cohttp_lwt_unix

let usage_msg = "guthub-user-activity <username>"
let username = ref ""
let get_username un = username := un
let () = Arg.parse [] get_username usage_msg
let url = "https://api.github.com/users/" ^ !username ^ "/events"
(* let () = Printf.printf "url: %s\n" url *)

let get_data_from_string s =
  match Yojson.Safe.from_string s with
  | `Assoc kv_list -> (
      let find key =
        match List.Assoc.find ~equal:String.equal kv_list key with
        | None | Some (`String "") -> None
        | Some s -> Some (Yojson.Safe.to_string s)
      in
      match find "Abstract" with Some _ as x -> x | None -> find "Definition")
  | _ -> None

let body =
  let url = Uri.of_string url in
  Client.get url >>= fun (resp, body) ->
  let code = resp |> Response.status |> Code.code_of_status in
  match code with
  | 200 ->
      Cohttp_lwt.Body.to_string body >>= fun body ->
      (* let json = Yojson.Safe.from_string body in *)
      let a = get_data_from_string body in
      Lwt.return (Ok a)
  | code ->
      let err =
        code |> Int.to_string |> Printf.sprintf "error with status code %s"
      in
      Lwt.return (Error err)
(* Printf.printf "Response code: %d\n" code; *)
(* Printf.printf "Received a response from the server"; *)
(* Printf.printf "Headers: %s\n" (resp |> Response.headers |> Header.to_string); *)
(* body |> Cohttp_lwt.Body.to_string >|= fun body -> *)
(* Printf.printf "Body of length: %d\n" (String.length body); *)
(* body *)
(* let json = Yojson.Safe.from_string body in *)

let () =
  let body = Lwt_main.run body in
  match body with Ok _ -> print_endline "OK" | Error msg -> print_endline msg

(* body |> String.length |> Int.to_string |> Printf.printf "length: %s\n" *)
