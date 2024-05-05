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

class ChatDetailPage extends StatefulWidget{
  final String initialChatId;

  ChatDetailPage({Key? key, required this.initialChatId}) : super(key: key);

  @override
  _ChatDetailPageState createState() => _ChatDetailPageState();
}


class _ChatDetailPageState extends State<ChatDetailPage> {
  late String chatId;
  List<local.ChatMessage> _messages = [];

  TextEditingController _messageController = TextEditingController();

  @override
  void initState() {
    super.initState();
    chatId = widget.initialChatId;
    _loadMessages();
  }

  void _loadMessages() async {
    if (chatId.isEmpty) {
      return;
    }

    final apiUrl = 'http://localhost:8080/chat_messages/${chatId}';
    final uri = Uri.parse(apiUrl).replace(queryParameters: {
     'limit': '100', // todo adjust as needed
      'timestamp': DateTime.now().subtract(Duration(days: 7)).millisecondsSinceEpoch.toString(),
    });

    final authService = Provider.of<AuthService>(context, listen: false);

    String? authToken = await authService.getCurrentUserIdToken();

    try {
      print('Fetching messages from $uri');
      final response = await http.get(
          uri,
          headers: {
            'Authorization': 'Bearer $authToken', // Include the authorization header
            'Content-Type': 'application/json',
          },
      ).timeout(Duration(seconds: 10));
      if (response.statusCode == 200) {
        final messagesResponse = GetChatMessagesResponse.fromJson(jsonDecode(response.body) as Map<String, dynamic>);

        setState(() {
          _messages = messagesResponse.Messages.map((message) => local.ChatMessage(
            IsFromCurrentUser: message.IsFromCurrentUser,
            Text: message.Text,
            Timestamp: message.Timestamp,
          )).toList();

          _messages.forEach((message) {
            print('Message: ${message.Text}');
            print('Timestamp: ${message.Timestamp}');
            print('IsFromCurrentUser: ${message.IsFromCurrentUser}');
          });
        });
      } else {
       // print('Server error: ${response.body}');
        print('Failed to fetch messages: ${response.statusCode}');
      }
    } catch (e) {
      print('Error fetching messages: $e');
    }
  }

  // send post request to server for adding message
  Future<SendChatMessageResponse> _sendMessage(SendChatMessageRequest message) async {
    const apiUrl = 'http://localhost:8080/chat_messages';
    final uri = Uri.parse(apiUrl);

    final authService = Provider.of<AuthService>(context, listen: false);

    String? authToken = await authService.getCurrentUserIdToken();

    try {
      print('Sending message to $uri');
      final response = await http.post(
        uri,
        headers: {
          'Authorization': 'Bearer $authToken', // Include the authorization header
          'Content-Type': 'application/json',
        },
        body: jsonEncode(message.toJson()),
      ).timeout(Duration(seconds: 10));
      if (response.statusCode == 200) {
        final responseMessage = SendChatMessageResponse.fromJson(jsonDecode(response.body));

        return responseMessage;
      } else {
        print('Server error: ${response.body}');
        throw Exception('Failed to send message');
      }
    } catch (e) {
      print('Error sending message: $e');
      throw Exception('Failed to send message');
    }
  }

  void _addMessage(local.ChatMessage message) {
    setState(() {
      _messages.add(message);
    });
  }

  void _setChatId(String newChatId) {
    setState(() {
      chatId = newChatId;
    });
  }

  void _handleSendPressed(String text) {
    var timestamp = DateTime.now().millisecondsSinceEpoch;

    _addMessage(local.ChatMessage(IsFromCurrentUser: true, Text: text, Timestamp: timestamp));

    final message = SendChatMessageRequest(
      ChatId: chatId,
      Text: text,
      Timestamp: timestamp,
    );

    // Send the message to the server
    _sendMessage(message).then((responseMessage) {
      if (responseMessage.CreatedChat != null) {
        Provider.of<localChat.ChatModel>(context, listen: false).addChat(
          localChat.Chat(
            ChatId: responseMessage.CreatedChat!.ChatId,
            Timestamp: responseMessage.CreatedChat!.Timestamp,
            Title: responseMessage.CreatedChat!.Title,
          ),
        );

        _setChatId(responseMessage.CreatedChat!.ChatId);
      }

      _addMessage(
        local.ChatMessage(
          IsFromCurrentUser: false,
          Text: responseMessage.Reply,
          Timestamp: responseMessage.Timestamp,
        ),
      );
    }).catchError((e) {
      print('Error sending message: $e');
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Stack(
        children: <Widget>[
          ListView.builder(
            itemCount: _messages.length,
            shrinkWrap: true,
            padding: EdgeInsets.only(top: 10,bottom: 10),
            physics: NeverScrollableScrollPhysics(),
            itemBuilder: (context, index){
              return Container(
                padding: EdgeInsets.only(left: 14,right: 14,top: 10,bottom: 10),
                child: Align(
                  alignment: (_messages[index].IsFromCurrentUser?Alignment.topRight:Alignment.topLeft),
                  child: Container(
                    decoration: BoxDecoration(
                      borderRadius: BorderRadius.circular(20),
                      color: (_messages[index].IsFromCurrentUser?Colors.blue[200]:Colors.grey.shade200),
                    ),
                    padding: EdgeInsets.all(16),
                    child: Text(_messages[index].Text, style: TextStyle(fontSize: 15),),
                  ),
                ),
              );
            },
          ),
          Align(
            alignment: Alignment.bottomLeft,
            child: Container(
              padding: EdgeInsets.only(left: 10,bottom: 10,top: 10),
              height: 60,
              width: double.infinity,
              color: Colors.white,
              child: Row(
                children: <Widget>[
                  GestureDetector(
                    onTap: (){
                    },
                    child: Container(
                      height: 30,
                      width: 30,
                      decoration: BoxDecoration(
                        color: Colors.lightBlue,
                        borderRadius: BorderRadius.circular(30),
                      ),
                      child: Icon(Icons.add, color: Colors.white, size: 20, ),
                    ),
                  ),
                  SizedBox(width: 15,),
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
                  SizedBox(width: 15,),
                  FloatingActionButton(
                    onPressed: (){
                      if (_messageController.text.trim().isNotEmpty) {
                        _handleSendPressed(_messageController.text.trim());
                        _messageController.clear();
                      }
                    },
                    child: Icon(Icons.send,color: Colors.white,size: 18,),
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
