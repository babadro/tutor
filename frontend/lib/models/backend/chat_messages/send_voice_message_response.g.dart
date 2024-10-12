// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'send_voice_message_response.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

SendVoiceMessageResponse _$SendVoiceMessageResponseFromJson(
        Map<String, dynamic> json) =>
    SendVoiceMessageResponse(
      UserText: json['usrTxt'] as String,
      UserAudioURL: json['usrAudio'] as String,
      UserMessageTime: (json['usrTime'] as num).toInt(),
      ReplyText: json['replyTxt'] as String? ?? '',
      ReplyAudioURL: json['replyAudio'] as String? ?? '',
      ReplyTime: (json['replyTime'] as num?)?.toInt() ?? 0,
      CreatedChat: json['chat'] == null
          ? null
          : Chat.fromJson(json['chat'] as Map<String, dynamic>),
    );

Map<String, dynamic> _$SendVoiceMessageResponseToJson(
        SendVoiceMessageResponse instance) =>
    <String, dynamic>{
      'usrTxt': instance.UserText,
      'usrAudio': instance.UserAudioURL,
      'usrTime': instance.UserMessageTime,
      'replyTxt': instance.ReplyText,
      'replyAudio': instance.ReplyAudioURL,
      'replyTime': instance.ReplyTime,
      'chat': instance.CreatedChat,
    };
