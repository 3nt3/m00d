import 'package:flutter/material.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:google_sign_in/google_sign_in.dart';

import '../api.dart';

class MyLoginPage extends StatefulWidget {
  const MyLoginPage({Key? key}) : super(key: key);

  @override
  State<StatefulWidget> createState() => _MyLoginPageState();
}

class _MyLoginPageState extends State<MyLoginPage> {
  @override
  void initState() {
    // TODO: implement initState
    super.initState();
    debugPrint('this is the route were on');
  }

  Future<void> _onPress() async {
    late final String? token;
    try {
      token = await initAuth(visible: true);
    } catch (e) {
      if (!mounted) return;
      _handleError(context, e.toString());
    }

    // I'm not sure if this is ever actually called since login() should throw an exception whenever something goes wrong
    if (token == null) {
      if (!mounted) return; // what does this even do?
      _handleError(context, "something bad happened");
    }
  }

  _handleError(BuildContext context, String error) {
    showDialog(
      context: context,
      builder: (BuildContext builder) => Expanded(
        child: AlertDialog(
          title: const Text("Login failed"),
          content: Text(error),
        ),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Container(
        color: Colors.black,
        child: Padding(
          padding: const EdgeInsets.symmetric(vertical: 50, horizontal: 20),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              const Text('Log in to m00d',
                  style: TextStyle(color: Colors.white, fontSize: 30)),
              const SizedBox(height: 20),
              ElevatedButton(
                style: ElevatedButton.styleFrom(
                    primary: Colors.grey.shade800,
                    padding: const EdgeInsets.all(20)),
                onPressed: _onPress,
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    SizedBox(
                        height: 30,
                        child: Image.asset('assets/google_logo.png')),
                    const SizedBox(width: 10),
                    const Text(
                      'GOOGLE',
                      style: TextStyle(
                        fontSize: 20,
                      ),
                    ),
                  ],
                ),
              )
            ],
          ),
        ),
      ),
    );
  }
}
