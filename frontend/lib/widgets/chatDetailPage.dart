import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:tutor/models/backend/chat_messages/send_text_message_request.dart';
import 'package:tutor/models/local/chat/chats.dart' as localChat;
import 'package:tutor/services/audio_player_service.dart';
import 'package:tutor/services/auth_service.dart';
import 'package:tutor/models/local/chat/chat_message.dart' as local;
import 'package:tutor/services/chat_service.dart';
import 'package:flutter_sound_platform_interface/flutter_sound_recorder_platform_interface.dart';

import 'package:tutor/services/audio_recorder_service.dart';

import 'message.dart';

typedef _Fn = void Function();
const theSource = AudioSource.microphone;

class ChatDetailPage extends StatefulWidget {
  final String initialChatId;
  final AudioRecorderService mRecorder;

  ChatDetailPage({Key? key, required this.initialChatId, required this.mRecorder}) : super(key: key);

  @override
  _ChatDetailPageState createState() => _ChatDetailPageState();
}

class _ChatDetailPageState extends State<ChatDetailPage> {
  late String chatId;
  late ChatService _chatService;
  List<local.ChatMessage> _messages = [];
  final ScrollController _scrollController = ScrollController();

  TextEditingController _messageController = TextEditingController();
  AudioRecorderService get _mRecorder => widget.mRecorder;

  bool _isRecording = false;
  bool _isSending = false;

  @override
  void initState() {
    chatId = widget.initialChatId;
    _chatService = ChatService(Provider.of<AuthService>(context, listen: false));
    _loadMessages();
    _mRecorder.init();

    super.initState();
  }

  @override
  void dispose() {
    _scrollController.dispose();
    super.dispose();
  }

  void _loadMessages() async {
    var loadMessagesResult = await _chatService.loadMessages(chatId);
    if (!loadMessagesResult.success) {
      print('Failed to load messages: ${loadMessagesResult.errorMessage}');
      // todo: show error message
      return;
    }

    setState(() {
      _messages = loadMessagesResult.data!;
      _scrollToBottom();
    });
  }

  void _addMessage(local.ChatMessage message) {
    setState(() {
      _messages.add(message);
      _scrollToBottom();
    });
  }

  void _handleSendPressed(String text) async {
    var timestamp = DateTime.now().millisecondsSinceEpoch;
    var message = SendTextMessageRequest(
      ChatId: chatId,
      Text: text,
      Timestamp: timestamp,
    );

    _addMessage(
        local.ChatMessage(
          IsFromCurrentUser: true,
          Text: text,
          Timestamp: timestamp,
        )
    );

    var sendResult = await _chatService.sendMessage(message);
    if (!sendResult.success) {
      print('Failed to send message: ${sendResult.errorMessage}');
      // todo: show error message
      return;
    }

    if (sendResult.data!.createdChat != '') {
      switchToNewChat(sendResult.data!.createdChat);
    };

    _addMessage(sendResult.data!.message);
  }

  void switchToNewChat(localChat.Chat createdChat) {
    Provider.of<localChat.ChatModel>(context, listen: false).addChat(
      createdChat,
    );

    setState(() {
      chatId = createdChat.ChatId;
    });
  }

  void record() {
    _mRecorder.record(
    ).then((_) {
      setState(() {
        _isRecording = true;
      });
    });
  }

  void stopRecorder() async {
    setState(() {
      _isRecording = false;
      _isSending = true;
    });

    await _mRecorder.stopRecording().then((value) {
      _chatService.sendVoiceMessage(value ?? '', chatId).then((value) {
        setState(() {
          _isSending = false;
        });

        if (!value.success) {
          print('Failed to send voice message: ${value.errorMessage}');
        } else {
          final res = value.data!;

          if (value.data!.createdChat.ChatId != '') {
            switchToNewChat(value.data!.createdChat);
          };

          _addMessage(res.userMessage);
          _addMessage(res.replyMessage);
        }
      });
    });
  }

  _Fn? getRecorderFn() {
    if (!_mRecorder.inited) {
      return null;
    }
    return _mRecorder.isStopped ? record : stopRecorder;
  }

  void _cancelRecording() async {
    await _mRecorder.stopRecording();
    setState(() {
      _isRecording = false;
    });
  }

  void _scrollToBottom() {
    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (_scrollController.hasClients) {
        _scrollController.jumpTo(_scrollController.position.maxScrollExtent);
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    return ChangeNotifierProvider(
        create: (context) => AudioPlayerService(),
        child: Scaffold(
          body: Stack(
            children: <Widget>[
              Column(
                children: <Widget>[
                  Expanded(
                    child: ListView.builder(
                      controller: _scrollController,
                      itemCount: _messages.length + (_isRecording || _isSending ? 1 : 0),
                      padding: EdgeInsets.only(top: 10, bottom: 70),
                      itemBuilder: (context, index) {
                        if (index < _messages.length) {
                          return MessageWidget(
                            key: ValueKey(_messages[index].Timestamp),
                            message: _messages[index],
                          );
                        }

                        return Align(
                          alignment: Alignment.topRight,
                          child: Padding(
                            padding: const EdgeInsets.all(20),
                            child: Container(
                              decoration: BoxDecoration(
                                borderRadius: BorderRadius.circular(20),
                                color: _isRecording ? Colors.orange : Colors.blue,
                              ),
                              padding: EdgeInsets.all(16),
                              child: Row(
                                mainAxisSize: MainAxisSize.min,
                                children: [
                                  Icon(_isRecording ? Icons.mic : Icons.send, color: Colors.white),
                                  SizedBox(width: 8),
                                  Text(
                                    _isRecording ? 'Recording' : 'Sending',
                                    style: TextStyle(color: Colors.white),
                                  ),
                                  SizedBox(width: 8),
                                  CircularProgressIndicator(),
                                ],
                              ),
                            ),
                          ),
                        );
                      },
                    ),
                  ),
                ],
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
                        onTap: getRecorderFn(),
                        child: MouseRegion(
                          cursor: SystemMouseCursors.click,
                          child: Container(
                            height: 30,
                            width: 30,
                            decoration: BoxDecoration(
                              color: _mRecorder.isRecording ? Colors.red : Colors.lightBlue,
                              borderRadius: BorderRadius.circular(30),
                            ),
                            child: Icon(_mRecorder.isRecording ? Icons.stop : Icons.mic, color: Colors.white, size: 20),
                          ),
                        ),
                      ),
                      Visibility(
                        visible: _mRecorder.isRecording,
                        child: GestureDetector(
                          onTap: _cancelRecording,
                          child: MouseRegion(
                            cursor: SystemMouseCursors.click,
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
        )
    );
  }
}

