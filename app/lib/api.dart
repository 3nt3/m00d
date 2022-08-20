import 'dart:convert';
import 'dart:io';

import 'package:flutter/material.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:google_sign_in/google_sign_in.dart';

import 'models.dart';

import 'package:http/http.dart' as http;

const apiUrl = "http://192.168.178.196:8080";

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

Future<String?> initAuth({required bool visible}) async {
  const storage = FlutterSecureStorage();

  var jwt = await storage.read(key: "jwt");
  if (jwt != null) {
    debugPrint("token exists, attempting refresh");
    jwt = await _refreshToken(jwt, visible: visible);
  } else {
    debugPrint("no token stored, logging in");
    jwt = await login(visible: visible);
  }

  debugPrint("writing $jwt to storage");
  await storage.write(key: "jwt", value: jwt);
  return jwt;
}

Future<String?> _refreshToken(String jwt, {required bool visible}) async {
  var response = await http.post(Uri.parse("$apiUrl/refresh-token"),
      headers: {"Authorization": "Bearer $jwt"});

  if (response.statusCode >= 400) {
    try {
      return await login(visible: visible);
    } catch (_) {
      return null;
    }
  }

  return jsonDecode(response.body)['token'];
}

Future<String?> login({required bool visible}) async {
  // first try signing in silently, and then don't if that doesn't work
  debugPrint('refresh didn\'t work for some reason, logging in normally');

  final GoogleSignIn googleSignIn = GoogleSignIn(
      scopes: <String>["email"],
      serverClientId:
          '82145806916-vocueu5na49d2lgusnotbrjdd7ne77mp.apps.googleusercontent.com');
  var acc = await googleSignIn.signInSilently() ??
      (visible ? await googleSignIn.signIn() : null);
  if (acc == null) return null;

  final GoogleSignInAuthentication auth = await acc.authentication;
  var response = await http.post(Uri.parse("$apiUrl/login"),
      body: jsonEncode({"id_token": auth.idToken}),
      headers: <String, String>{
        HttpHeaders.contentTypeHeader: 'application/json'
      });

  return jsonDecode(response.body)['token'];
}
