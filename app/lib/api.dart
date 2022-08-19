<<<<<<< Updated upstream
import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:google_sign_in/google_sign_in.dart';

import 'models.dart';

import 'package:http/http.dart' as http;

const apiUrl = "http://gott-himalayas.nord:8080";

Future<List<Mood>> getMoods({required String token}) async {
  var resp = await http.get(Uri.parse('$apiUrl/moods'),
      headers: {'Authorization': 'Bearer $token'});

  if (resp.statusCode == 200) {
    return (jsonDecode(resp.body) as List<Map<String, dynamic>>)
        .map((map) => Mood.fromJson(map))
        .toList();
  } else {
    throw Exception(_extractError(resp));
  }
}

Future<User?> getUser({required String? token}) async {
  if (token == null) return null;

  var resp = await http.get(Uri.parse('$apiUrl/user'),
      headers: {'authorization': 'Bearer $token'});
  if (resp.statusCode == 200) {
    return User.fromJson(jsonDecode(resp.body));
  } else {
    throw Exception(_extractError(resp));
  }
}

String _extractError(http.Response resp) {
  Map<String, dynamic> body;
  try {
    body = jsonDecode(resp.body);
  } on FormatException {
    return 'Request failed with status ${resp.statusCode}. Body not parseable';
  }

  var error = body['error'] ?? '<No message>';

  return 'Request failed with status ${resp.statusCode}. Error: \'$error\'';
}

Future<String?> initAuth() async {
  final storage = const FlutterSecureStorage();

  var jwt = await storage.read(key: "jwt");
  if (jwt != null) {
    debugPrint("token exists, attempting refresh");
    jwt = await _refreshToken(jwt);
  } else {
    debugPrint("no token stored, logging in");
    jwt = await _login();
  }

  storage.write(key: "jwt", value: jwt);
  return jwt;
}

Future<String?> _refreshToken(String jwt) async {
  var response = await http.post(Uri.parse("$apiUrl/refresh-token"),
      headers: {"Authorization": "Bearer $jwt"});

  if (response.statusCode >= 400) {
    try {
      return await _login();
    } catch (_) {
      return null;
    }
  }

  return jsonDecode(response.body)['token'];
}

Future<String?> _login() async {
  // first try signing in silently, and then don't if that doesn't work
  debugPrint('refresh didn\'t work for some reason, logging in normally');

  final GoogleSignIn _googleSignIn = GoogleSignIn(
      scopes: <String>["email"],
      serverClientId:
          '82145806916-vocueu5na49d2lgusnotbrjdd7ne77mp.apps.googleusercontent.com');
  var acc = await googleSignIn.signInSilently() ?? await googleSignIn.signIn();
  if (acc == null) return null;

  final GoogleSignInAuthentication auth = await acc.authentication;
  debugPrint(auth.idToken);
  var response = await http.post(Uri.parse("$apiUrl/login"),
      body: jsonEncode({"id_token": auth.idToken}),
      headers: <String, String>{
        HttpHeaders.contentTypeHeader: 'application/json'
      });

  debugPrint(jsonDecode(response.body)['token']);

  return jsonDecode(response.body)['token'];
}

class _googleSignIn {}
=======
//final API_URL = "http://gott-himalayas.nord:8080";
final API_URL = "http://192.178.168.196:8080";
>>>>>>> Stashed changes
