import 'dart:typed_data';

import 'package:http/http.dart' as http;
import 'package:tutor/models/backend/chat_messages/send_text_message_request.dart';
import 'package:tutor/models/backend/chat_messages/send_text_message_response.dart';
import 'package:tutor/models/backend/chat_messages/get_chat_messages_response.dart';
import 'package:tutor/models/local/chat/chat_message.dart' as local;
import 'package:tutor/services/auth_service.dart';
import 'dart:convert';
import 'package:tutor/models/local/chat/chats.dart' as localChat;
import 'package:tutor/services/service_response.dart';
import 'package:tutor/models/backend/chats/get_chats_response.dart';
import 'package:tutor/models/backend/chat_messages/send_voice_message_response.dart';

class sendMessageResult {
  final local.ChatMessage message;
  final localChat.Chat createdChat;

  sendMessageResult(this.message, this.createdChat);
}

class sendVoiceMessageResult {
  final local.ChatMessage userMessage;
  final local.ChatMessage replyMessage;
  final localChat.Chat createdChat;

  sendVoiceMessageResult(this.userMessage, this.replyMessage, this.createdChat);
}

class ChatService {
  final AuthService _authService;

  ChatService(this._authService);

  Future<ServiceResult<List<local.ChatMessage>>> loadMessages(
      String chatId) async {
    if (chatId.isEmpty) {
      return ServiceResult.success([]);
    }

    final apiUrl = 'http://localhost:8080/chat_messages/$chatId';
    final uri = Uri.parse(apiUrl).replace(queryParameters: {
      'limit': '100'
      //'timestamp': DateTime.now().subtract(Duration(days: 356)).millisecondsSinceEpoch.toString(),
    });

    String? authToken = await _authService.getCurrentUserIdToken();

    try {
      final response = await http.get(uri, headers: {
        'Authorization': 'Bearer $authToken',
        'Content-Type': 'application/json',
      });

      if (response.statusCode == 200) {
        final decodedResponseBody = utf8.decode(response.bodyBytes);
        final messagesResponse =
            GetChatMessagesResponse.fromJson(jsonDecode(decodedResponseBody));
        return ServiceResult.success(
            messagesResponse.Messages.map((message) => local.ChatMessage(
                  IsFromCurrentUser: message.IsFromCurrentUser,
                  Text: message.Text,
                  Timestamp: message.Timestamp,
                  AudioUrl: message.AudioUrl,
                )).toList());
      } else {
        return ServiceResult.failure(
            errorMessage: 'Failed to fetch messages: ${response.statusCode}');
      }
    } catch (e) {
      return ServiceResult.failure(
          errorMessage: 'Failed to fetch messages: $e');
    }
  }

  Future<ServiceResult<sendMessageResult>> sendMessage(
      SendTextMessageRequest message) async {
    final apiUrl = 'http://localhost:8080/chat_messages';
    final uri = Uri.parse(apiUrl);
    String? authToken = await _authService.getCurrentUserIdToken();

    try {
      final response = await http.post(
        uri,
        headers: {
          'Authorization': 'Bearer $authToken',
          'Content-Type': 'application/json',
        },
        body: jsonEncode(message.toJson()),
      );
      if (response.statusCode == 200) {
        final resp =
            SendTextMessageResponse.fromJson(jsonDecode(response.body));

        return ServiceResult.success(
          sendMessageResult(
            local.ChatMessage(
              IsFromCurrentUser: false,
              Text: resp.Reply,
              Timestamp: resp.Timestamp,
            ),
            resp.CreatedChat != null
                ? localChat.Chat.fromChatResponse(resp.CreatedChat!)
                : localChat.Chat(
                    ChatId: '',
                    Timestamp: 0,
                    Title: '',
                    Type: localChat.ChatType.General),
          ),
        );
      } else {
        return ServiceResult.failure(
            errorMessage: 'Failed to send message: ${response.statusCode}');
      }
    } catch (e) {
      return ServiceResult.failure(errorMessage: 'Failed to send message: $e');
    }
  }

  Future<ServiceResult<sendVoiceMessageResult>> sendVoiceMessage(
      String audioFilePath, String chatId) async {
    final apiUrl = 'http://localhost:8080/chat_voice_messages';
    final uri = Uri.parse(apiUrl);
    String? authToken = await _authService.getCurrentUserIdToken();

    // todo will not work on mobile app
    Uint8List fileBytes = await http.readBytes(Uri.parse(audioFilePath));

    print('File size: ${fileBytes.length} bytes');

    try {
      final timestamp = DateTime.now();

      final request = http.MultipartRequest('POST', uri)
        ..headers['Authorization'] = 'Bearer $authToken'
        ..fields['chatId'] = chatId
        ..fields['timestamp'] = timestamp.millisecondsSinceEpoch.toString()
        ..files.add(
          await http.MultipartFile.fromBytes(
            'file',
            fileBytes,
            filename: 'audio_${timestamp.millisecondsSinceEpoch}.m4a',
          ),
        );

      final streamedResponse = await request.send();
      final response = await http.Response.fromStream(streamedResponse);

      if (response.statusCode == 200) {
        final resp =
            SendVoiceMessageResponse.fromJson(jsonDecode(response.body));

        return ServiceResult.success(
          sendVoiceMessageResult(
            local.ChatMessage(
              IsFromCurrentUser: true,
              Text: resp.UserText,
              AudioUrl: resp.UserAudioURL,
              Timestamp: resp.UserMessageTime,
            ),
            local.ChatMessage(
              IsFromCurrentUser: false,
              Text: resp.ReplyText,
              AudioUrl: resp.ReplyAudioURL,
              Timestamp: resp.ReplyTime,
            ),
            resp.CreatedChat != null
                ? localChat.Chat.fromChatResponse(resp.CreatedChat!)
                : localChat.Chat(
                    ChatId: '',
                    Timestamp: 0,
                    Title: '',
                    Type: localChat.ChatType.General),
          ),
        );
      } else {
        return ServiceResult.failure(
            errorMessage:
                'Failed to send voice message: ${response.statusCode}');
      }
    } catch (e) {
      return ServiceResult.failure(
          errorMessage: 'Failed to send voice message: $e');
    }
  }

  Future<ServiceResult<List<localChat.Chat>>> getChats() async {
    const apiUrl = 'http://localhost:8080/chats?limit=100&timestamp=0';
    final uri = Uri.parse(apiUrl);

    String? authToken = await _authService.getCurrentUserIdToken();

    try {
      print('Fetching chats from $uri');
      final response = await http.get(
        uri,
        headers: {
          'Authorization':
              'Bearer $authToken', // Include the authorization header
          'Content-Type': 'application/json',
        },
      ).timeout(Duration(seconds: 10));
      if (response.statusCode == 200) {
        final chatsResponse = GetChatsResponse.fromJson(
            jsonDecode(response.body) as Map<String, dynamic>);

        var chats =
            chatsResponse.Chats.map((e) => (localChat.Chat.fromChatResponse(e)))
                .toList();

        return ServiceResult.success(chats);
      } else {
        return ServiceResult.failure(
            errorMessage: 'Failed to fetch chats: ${response.statusCode}');
      }
    } catch (e) {
      return ServiceResult.failure(errorMessage: 'Failed to fetch chats: $e');
    }
  }
/*
  // Create chat
  Future<ServiceResult<localChat.Chat>> createChat(localChat.ChatType type) async {
    final apiUrl = 'http://localhost:8080/chats';
    final uri = Uri.parse(apiUrl);
    String? authToken = await _authService.getCurrentUserIdToken();

    try {
      final response = await http.post(
        uri,
        headers: {
          'Authorization': 'Bearer $authToken',
          'Content-Type': 'application/json',
        },
        body: jsonEncode(localChat.Chat.toChatRequest(type).toJson(),
      );

      if (response.statusCode == 200) {
        final chatResponse = localChat.Chat.fromChatResponse(
            jsonDecode(response.body) as Map<String, dynamic>);

        return ServiceResult.success(chatResponse);
      } else {
        return ServiceResult.failure(
            errorMessage: 'Failed to create chat: ${response.statusCode}');
      }
    } catch (e) {
      return ServiceResult.failure(errorMessage: 'Failed to create chat: $e');
    }
  }

 */
}
