import 'package:json_annotation/json_annotation.dart';

part 'chat.g.dart';

@JsonSerializable()
class Chat {
  @JsonKey(name: 'id')
  final String Id;

  @JsonKey(name: 'title')
  final String Title;

  Chat({
    required this.Id,
    required this.Title,
  });

  factory Chat.fromJson(Map<String, dynamic> json) => _$ChatFromJson(json);
  Map<String, dynamic> toJson() => _$ChatToJson(this);
}