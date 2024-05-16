import 'package:json_annotation/json_annotation.dart';

import 'package:tutor/models/backend/chats/get_chats_response.dart';

part 'send_voice_message_response.g.dart';

@JsonSerializable()
class SendVoiceMessageResponse {
  @JsonKey(name: 'usrTxt')
  final String UserText;

  @JsonKey(name: 'usrAudio')
  final String UserAudioURL;

  @JsonKey(name: 'usrTime')
  final int UserMessageTime;

  @JsonKey(name: 'replyTxt')
  final String ReplyText;

  @JsonKey(name: 'replyAudio')
  final String ReplyAudioURL;

  @JsonKey(name: 'replyTime')
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