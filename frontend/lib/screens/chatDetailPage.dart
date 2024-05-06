import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:tutor/models/backend/chat_messages/send_chat_message_request.dart';
import 'package:tutor/models/backend/chat_messages/send_chat_message_response.dart';
import 'package:tutor/models/local/chat/chats.dart' as localChat;
import 'package:tutor/services/auth_service.dart';
import 'package:tutor/models/backend/chat_messages/get_chat_messages_response.dart';
import 'package:tutor/models/local/chat/chat_message.dart' as local;
import 'package:flutter_sound/flutter_sound.dart';
import 'package:permission_handler/permission_handler.dart';

import '../services/chat_service.dart';

class ChatDetailPage extends StatefulWidget{
  final String initialChatId;

  ChatDetailPage({Key? key, required this.initialChatId}) : super(key: key);

  @override
  _ChatDetailPageState createState() => _ChatDetailPageState();
}


class _ChatDetailPageState extends State<ChatDetailPage> {
  late String chatId;
  late ChatService _chatService;
  List<local.ChatMessage> _messages = [];

  TextEditingController _messageController = TextEditingController();
  bool _isRecording = false;

  @override
  void initState() {
    super.initState();
    chatId = widget.initialChatId;
    _chatService = ChatService(context);
    _loadMessages();
  }

  void _loadMessages() async {
    var messages = await _chatService.loadMessages(chatId);
    setState(() {
      _messages = messages;
    });
  }

  void _addMessage(local.ChatMessage message) {
    setState(() {
      _messages.add(message);
    });
  }

  void _handleSendPressed(String text) async {
    var timestamp = DateTime.now().millisecondsSinceEpoch;
    var message = SendChatMessageRequest(
      ChatId: chatId,
      Text: text,
      Timestamp: timestamp,
    );

    try {
      var responseMessage = await _chatService.sendMessage(message);
      _addMessage(local.ChatMessage(
        IsFromCurrentUser: true,
        Text: text,
        Timestamp: timestamp,
      ));

      if (responseMessage.CreatedChat != null) {
        Provider.of<localChat.ChatModel>(context, listen: false).addChat(
          localChat.Chat(
            ChatId: responseMessage.CreatedChat!.ChatId,
            Timestamp: responseMessage.CreatedChat!.Timestamp,
            Title: responseMessage.CreatedChat!.Title,
          ),
        );

        setState(() {
          chatId = responseMessage.CreatedChat!.ChatId;
        });
      }

      _addMessage(local.ChatMessage(
        IsFromCurrentUser: false,
        Text: responseMessage.Reply,
        Timestamp: responseMessage.Timestamp,
      ));
    } catch (e) {
      print('Error sending message: $e');
    }
  }

  Future<void> _startRecording() async {
    // todo
  }

  Future<void> _stopAndSendRecording() async {
    // todo
  }

  void _cancelRecording() async {
    // todo
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Stack(
        children: <Widget>[
          ListView.builder(
            itemCount: _messages.length,
            shrinkWrap: true,
            padding: EdgeInsets.only(top: 10, bottom: 10),
            physics: NeverScrollableScrollPhysics(),
            itemBuilder: (context, index) {
              return Container(
                padding: EdgeInsets.only(left: 14, right: 14, top: 10, bottom: 10),
                child: Align(
                  alignment: (_messages[index].IsFromCurrentUser ? Alignment.topRight : Alignment.topLeft),
                  child: Container(
                    decoration: BoxDecoration(
                      borderRadius: BorderRadius.circular(20),
                      color: (_messages[index].IsFromCurrentUser ? Colors.blue[200] : Colors.grey.shade200),
                    ),
                    padding: EdgeInsets.all(16),
                    child: Text(_messages[index].Text, style: TextStyle(fontSize: 15)),
                  ),
                ),
              );
            },
          ),
          Align(
            alignment: Alignment.bottomLeft,
            child: Container(
              padding: EdgeInsets.only(left: 10, bottom: 10, top: 10),
              height: 60,
              width: double.infinity,
              color: Colors.white,
              child: Row(
                children: <Widget>[
                  GestureDetector(
                    onTap: () {
                      if (!_isRecording) {
                        _startRecording();
                      } else {
                        _stopAndSendRecording();
                      }
                    },
                    child: Container(
                      height: 30,
                      width: 30,
                      decoration: BoxDecoration(
                        color: _isRecording ? Colors.red : Colors.lightBlue,
                        borderRadius: BorderRadius.circular(30),
                      ),
                      child: Icon(_isRecording ? Icons.stop : Icons.mic, color: Colors.white, size: 20),
                    ),
                  ),
                  Visibility(
                    visible: _isRecording,
                    child: GestureDetector(
                      onTap: _cancelRecording,
                      child: Container(
                        height: 30,
                        width: 30,
                        decoration: BoxDecoration(
                          color: Colors.black,
                          borderRadius: BorderRadius.circular(30),
                        ),
                        child: Icon(Icons.delete, color: Colors.white, size: 20),
                      ),
                    ),
                  ),
                  SizedBox(width: 15),
                  Expanded(
                    child: TextField(
                      decoration: InputDecoration(
                          hintText: "Write message...",
                          hintStyle: TextStyle(color: Colors.black54),
                          border: InputBorder.none
                      ),
                      controller: _messageController,
                    ),
                  ),
                  SizedBox(width: 15),
                  FloatingActionButton(
                    onPressed: () {
                      if (_messageController.text.trim().isNotEmpty) {
                        _handleSendPressed(_messageController.text.trim());
                        _messageController.clear();
                      }
                    },
                    child: Icon(Icons.send, color: Colors.white, size: 18),
                    backgroundColor: Colors.blue,
                    elevation: 0,
                  ),
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }
}
