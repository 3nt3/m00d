import '../models.dart';
import 'package:flutter/material.dart';

class UserProvider with ChangeNotifier {
  User? _user;
  String? _token;
  bool _loading = true;

  User? get user => _user;
  String? get token => _token;
  bool get loading => _loading;

  UserProvider() {}

  set user(User? user) {
    _user = user;
    notifyListeners();
  }

  set token(String? token) {
    _token = token;
    notifyListeners();
  }

  set loading(bool loading) {
    _loading = loading;
    notifyListeners();
  }
}
