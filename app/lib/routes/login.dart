import 'package:flutter/material.dart';

class MyLoginPage extends StatefulWidget {
  const MyLoginPage({Key? key}) : super(key: key);

  @override
  State<StatefulWidget> createState() {
    return _MyLoginPageState();
  }
}

class _MyLoginPageState extends State<MyLoginPage> {
  final GoogleSignIn _googleSignIn = GoogleSignIn(
      scopes: <String>["email"],
      serverClientId:
          '82145806916-vocueu5na49d2lgusnotbrjdd7ne77mp.apps.googleusercontent.com');

  final _storage = const FlutterSecureStorage();
  String _token = '';

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Center(
        child: Column(
          children: [
            const Text('Login'),
            ElevatedButton(
              // FIXME
              onPressed: () {},
              child: Row(
                children: [
                  Image.asset('google_logo.png'),
                  const Text('Google'),
                ],
              ),
            )
          ],
        ),
      ),
    );
  }
}
