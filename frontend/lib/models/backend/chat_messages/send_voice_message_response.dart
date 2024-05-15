import 'package:json_annotation/json_annotation.dart';

import 'package:tutor/models/backend/chats/get_chats_response.dart';

part 'send_voice_message_response.g.dart';

@JsonSerializable()
class SendVoiceMessageResponse {
  @JsonKey(name: 'usr_txt')
  final String UserText;

  @JsonKey(name: 'usr_audio')
  final String UserAudioURL;

  @JsonKey(name: 'usr_time')
  final int UserMessageTime;

  @JsonKey(name: 'reply_txt')
  final String ReplyText;

  @JsonKey(name: 'reply_audio')
  final String ReplyAudioURL;

  @JsonKey(name: 'reply_time')
  final int ReplyTime;

  @JsonKey(name: 'chat')
  final Chat? CreatedChat;

  SendVoiceMessageResponse({
    required this.UserText,
    required this.UserAudioURL,
    required this.UserMessageTime,
    required this.ReplyText,
    required this.ReplyAudioURL,
    required this.ReplyTime,
    this.CreatedChat,
  });

  factory SendVoiceMessageResponse.fromJson(Map<String, dynamic> json) => _$SendVoiceMessageResponseFromJson(json);
  Map<String, dynamic> toJson() => _$SendVoiceMessageResponseToJson(this);
}