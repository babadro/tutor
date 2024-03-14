import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:provider/provider.dart';
import '../models/chat_message.dart';
import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:tutor/services/auth_service.dart';

class ChatDetailPage extends StatefulWidget{
  @override
  _ChatDetailPageState createState() => _ChatDetailPageState();
}


class _ChatDetailPageState extends State<ChatDetailPage> {
  List<ChatMessage> _messages = [];

  final AuthService _authService = AuthService();

  @override
  void initState() {
    super.initState();
    _loadMessages();
  }

  void _loadMessages() async {
    const chatId = 'your_chat_id_here'; // Use the appropriate chatId for your use case
    const apiUrl = 'http://localhost:8080/chat_messages/1234';
    final uri = Uri.parse(apiUrl).replace(queryParameters: {
     'limit': '10', // Example, adjust as needed
      // 'before_timestamp': 'your_timestamp_here', // Uncomment and adjust if needed
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
        final messages = (jsonDecode(response.body) as List)
            .map((e) => ChatMessage.fromJson(e as Map<String, dynamic>))
            .toList();

        setState(() {
          _messages = messages;
        });
      } else {
        print('Server error: ${response.body}');
      }
    } catch (e) {
      print('Error fetching messages: $e');
    }
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
                    ),
                  ),
                  SizedBox(width: 15,),
                  FloatingActionButton(
                    onPressed: (){},
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
