import 'package:http/http.dart' as http;
import 'package:tutor/models/backend/chat_messages/send_chat_message_request.dart';
import 'package:tutor/models/backend/chat_messages/send_chat_message_response.dart';
import 'package:tutor/models/backend/chat_messages/get_chat_messages_response.dart';
import 'package:tutor/models/local/chat/chat_message.dart' as local;
import 'package:tutor/services/auth_service.dart';
import 'package:provider/provider.dart';
import 'package:flutter/material.dart';
import 'dart:convert';
import 'package:tutor/models/local/chat/chats.dart' as localChat;

import '../models/backend/chats/get_chats_response.dart';

class ChatService {
  final BuildContext context;

  ChatService(this.context);

  Future<List<local.ChatMessage>> loadMessages(String chatId) async {
    if (chatId.isEmpty) {
      return [];
    }

    final apiUrl = 'http://localhost:8080/chat_messages/$chatId';
    final uri = Uri.parse(apiUrl).replace(queryParameters: {
      'limit': '100',
      'timestamp': DateTime.now().subtract(Duration(days: 7)).millisecondsSinceEpoch.toString(),
    });

    final authService = Provider.of<AuthService>(context, listen: false);
    String? authToken = await authService.getCurrentUserIdToken();

    try {
      final response = await http.get(uri, headers: {
        'Authorization': 'Bearer $authToken',
        'Content-Type': 'application/json',
      });

      if (response.statusCode == 200) {
        final messagesResponse = GetChatMessagesResponse.fromJson(jsonDecode(response.body));
        return messagesResponse.Messages.map((message) => local.ChatMessage(
          IsFromCurrentUser: message.IsFromCurrentUser,
          Text: message.Text,
          Timestamp: message.Timestamp,
        )).toList();
      } else {
        print('Failed to fetch messages: ${response.statusCode}');
        return [];
      }
    } catch (e) {
      print('Error fetching messages: $e');
      return [];
    }
  }

  Future<SendChatMessageResponse> sendMessage(SendChatMessageRequest message) async {
    final apiUrl = 'http://localhost:8080/chat_messages';
    final uri = Uri.parse(apiUrl);
    final authService = Provider.of<AuthService>(context, listen: false);
    String? authToken = await authService.getCurrentUserIdToken();

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
        return SendChatMessageResponse.fromJson(jsonDecode(response.body));
      } else {
        print('Server error: ${response.body}');
        throw Exception('Failed to send message');
      }
    } catch (e) {
      print('Error sending message: $e');
      throw Exception('Failed to send message');
    }
  }

  Future<List<localChat.Chat>> getChats() async {
    const apiUrl = 'http://localhost:8080/chats?limit=100&timestamp=0';
    final uri = Uri.parse(apiUrl);

    final authService = Provider.of<AuthService>(context, listen: false);

    String? authToken = await authService.getCurrentUserIdToken();

    try {
      print('Fetching chats from $uri');
      final response = await http.get(
        uri,
        headers: {
          'Authorization': 'Bearer $authToken', // Include the authorization header
          'Content-Type': 'application/json',
        },
      ).timeout(Duration(seconds: 10));
      if (response.statusCode == 200) {
        final chatsResponse = GetChatsResponse.fromJson(jsonDecode(response.body) as Map<String, dynamic>);

        var chats = chatsResponse.Chats.map((e) => ( localChat.Chat(
          ChatId: e.ChatId,
          Timestamp: e.Timestamp,
          Title: e.Title,
        ))).toList();

        return chats;
      } else {
        print('Failed to fetch chats: ${response.statusCode}');
        throw Exception('Failed to fetch chats');
      }
    } catch (e) {
      print('Failed to fetch chats: $e');
      throw Exception('Failed to fetch chats');
    }
  }
}
