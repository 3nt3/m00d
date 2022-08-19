import 'package:json_annotation/json_annotation.dart';

part 'models.g.dart';

@JsonSerializable()
class Mood {
  final int id;
  final int mood;
  DateTime createdAt;

  Mood({required this.id, required this.mood, required this.createdAt});

  /// Connect the generated [_$MoodFromJson] function to the `fromJson`
  /// factory.
  factory Mood.fromJson(Map<String, dynamic> json) => _$MoodFromJson(json);

  /// Connect the generated [_$MoodToJson] function to the `toJson` method.
  Map<String, dynamic> toJson() => _$MoodToJson(this);
}

@JsonSerializable()
class User {
  final int id;
  final String email;

  User({required this.id, required this.email});

  /// Connect the generated [_$UserFromJson] function to the `fromJson`
  /// factory.
  factory User.fromJson(Map<String, dynamic> json) => _$UserFromJson(json);

  /// Connect the generated [_$UserToJson] function to the `toJson` method.
  Map<String, dynamic> toJson() => _$UserToJson(this);
}
