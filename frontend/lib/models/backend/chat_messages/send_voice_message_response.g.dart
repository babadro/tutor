// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'send_voice_message_response.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

SendVoiceMessageResponse _$SendVoiceMessageResponseFromJson(
        Map<String, dynamic> json) =>
    SendVoiceMessageResponse(
      UserText: json['usr_txt'] as String,
      UserAudioURL: json['usr_audio'] as String,
      UserMessageTime: (json['usr_time'] as num).toInt(),
      ReplyText: json['reply_txt'] as String,
      ReplyAudioURL: json['reply_audio'] as String,
      ReplyTime: (json['reply_time'] as num).toInt(),
      CreatedChat: json['chat'] == null
          ? null
          : Chat.fromJson(json['chat'] as Map<String, dynamic>),
    );

Map<String, dynamic> _$SendVoiceMessageResponseToJson(
        SendVoiceMessageResponse instance) =>
    <String, dynamic>{
      'usr_txt': instance.UserText,
      'usr_audio': instance.UserAudioURL,
      'usr_time': instance.UserMessageTime,
      'reply_txt': instance.ReplyText,
      'reply_audio': instance.ReplyAudioURL,
      'reply_time': instance.ReplyTime,
      'chat': instance.CreatedChat,
    };
