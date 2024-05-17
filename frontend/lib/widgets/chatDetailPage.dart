import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:tutor/models/backend/chat_messages/send_text_message_request.dart';
import 'package:tutor/models/local/chat/chats.dart' as localChat;
import 'package:tutor/services/auth_service.dart';
import 'package:tutor/models/local/chat/chat_message.dart' as local;
import 'package:tutor/services/chat_service.dart';
import 'dart:async';
import 'package:audio_session/audio_session.dart';
import 'package:flutter/foundation.dart' show kIsWeb;
import 'package:flutter_sound/flutter_sound.dart';
import 'package:flutter_sound_platform_interface/flutter_sound_recorder_platform_interface.dart';
import 'package:permission_handler/permission_handler.dart';

typedef _Fn = void Function();
const theSource = AudioSource.microphone;

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

  Codec _codec = Codec.aacMP4;
  String _mPath = 'tau_file.mp4';
  FlutterSoundPlayer? _mPlayer = FlutterSoundPlayer();
  FlutterSoundRecorder? _mRecorder = FlutterSoundRecorder();
  bool _mPlayerIsInited = false;
  bool _mRecorderIsInited = false;
  bool _mplaybackReady = false;

  @override
  void initState() {
    openTheRecorder().then((value) {
      setState(() {
        _mRecorderIsInited = true;
      });
    });

    _mPlayer!.openPlayer().then((value) {
      setState(() {
        _mPlayerIsInited = true;
      });
    });

    chatId = widget.initialChatId;
    _chatService = ChatService(Provider.of<AuthService>(context, listen: false));
    _loadMessages();

    super.initState();
  }

  @override
  void dispose() {
    _mPlayer!.closePlayer();
    _mPlayer = null;

    _mRecorder!.closeRecorder();
    _mRecorder = null;
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
    });
  }

  void _addMessage(local.ChatMessage message) {
    setState(() {
      _messages.add(message);
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

    var createdChat = sendResult.data!.createdChat;

    if (createdChat.ChatId != '') {
      Provider.of<localChat.ChatModel>(context, listen: false).addChat(
        createdChat,
      );

      setState(() {
        chatId = createdChat.ChatId;
      });
    }

    _addMessage(sendResult.data!.message);
  }

  Future<void> openTheRecorder() async {
    if (!kIsWeb) {
      var status = await Permission.microphone.request();
      if (status != PermissionStatus.granted) {
        throw RecordingPermissionException('Microphone permission not granted');
      }
    }
    await _mRecorder!.openRecorder();
    if (!await _mRecorder!.isEncoderSupported(_codec) && kIsWeb) {
      _codec = Codec.opusWebM;
      _mPath = 'tau_file.webm';
      if (!await _mRecorder!.isEncoderSupported(_codec) && kIsWeb) {
        _mRecorderIsInited = true;
        return;
      }
    }
    final session = await AudioSession.instance;
    await session.configure(AudioSessionConfiguration(
      avAudioSessionCategory: AVAudioSessionCategory.playAndRecord,
      avAudioSessionCategoryOptions:
      AVAudioSessionCategoryOptions.allowBluetooth |
      AVAudioSessionCategoryOptions.defaultToSpeaker,
      avAudioSessionMode: AVAudioSessionMode.spokenAudio,
      avAudioSessionRouteSharingPolicy:
      AVAudioSessionRouteSharingPolicy.defaultPolicy,
      avAudioSessionSetActiveOptions: AVAudioSessionSetActiveOptions.none,
      androidAudioAttributes: const AndroidAudioAttributes(
        contentType: AndroidAudioContentType.speech,
        flags: AndroidAudioFlags.none,
        usage: AndroidAudioUsage.voiceCommunication,
      ),
      androidAudioFocusGainType: AndroidAudioFocusGainType.gain,
      androidWillPauseWhenDucked: true,
    ));

    _mRecorderIsInited = true;
  }

  void record() {
    _mRecorder!
        .startRecorder(
      toFile: _mPath,
      codec: _codec,
      audioSource: theSource,
    )
        .then((value) {
      setState(() {});
    });
  }

  void stopRecorder() async {
    await _mRecorder!.stopRecorder().then((value) {
      // log value
      print('Value is: $value');

      _chatService.sendVoiceMessage(value ?? '', chatId).then((value) {
        if (!value.success) {
          print('Failed to send voice message: ${value.errorMessage}');
        } else {
          print('Voice message sent');
          final res = value.data!;

          _addMessage(res.userMessage);
          _addMessage(res.replyMessage);

          // print reply message
          print('Reply message: ${res.replyMessage.AudioUrl}');
        }
      });

      setState(() {
       // var url = value;
        _mplaybackReady = true;
      });

      var timestamp = DateTime.now().millisecondsSinceEpoch;

      _addMessage(
          local.ChatMessage(
            IsFromCurrentUser: true,
            Text: "audio with local path",
            Timestamp: timestamp,
            AudioUrl: value ?? '',
          )
      );
    });
  }

  void play(String path) {
    print('Playing $path');
    assert(_mPlayerIsInited &&
        _mplaybackReady &&
        _mRecorder!.isStopped &&
        _mPlayer!.isStopped);
    _mPlayer!
        .startPlayer(
        fromURI: path,
        //codec: kIsWeb ? Codec.opusWebM : Codec.aacADTS,
        whenFinished: () {
          setState(() {});
        })
        .then((value) {
      setState(() {});
    });
  }

  void stopPlayer() {
    _mPlayer!.stopPlayer().then((value) {
      setState(() {});
    });
  }

  _Fn? getRecorderFn() {
    if (!_mRecorderIsInited || !_mPlayer!.isStopped) {
      return null;
    }
    return _mRecorder!.isStopped ? record : stopRecorder;
  }

  _Fn? getPlaybackFn(String path) {
    if (!_mPlayerIsInited || !_mplaybackReady || !_mRecorder!.isStopped) {
      return null;
    }
    return _mPlayer!.isStopped ? (){play(path);} : stopPlayer;
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
                    child: Column(
                      children: <Widget>[
                        Text(_messages[index].Text, style: TextStyle(fontSize: 15)),
                        Visibility(
                          visible: _messages[index].AudioUrl != '',
                          child: GestureDetector(
                            onTap: getPlaybackFn(_mPath),
                            child: Container(
                              height: 30,
                              width: 30,
                              decoration: BoxDecoration(
                                color: Colors.black,
                                borderRadius: BorderRadius.circular(30),
                              ),
                              child: Icon(Icons.play_arrow, color: Colors.white, size: 20),
                            ),
                          ),
                        ),
                      ]
                    ),
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
                    onTap: getRecorderFn(),
                    child: Container(
                      height: 30,
                      width: 30,
                      decoration: BoxDecoration(
                        color: _mRecorder!.isRecording ? Colors.red : Colors.lightBlue,
                        borderRadius: BorderRadius.circular(30),
                      ),
                      child: Icon(_mRecorder!.isRecording ? Icons.stop : Icons.mic, color: Colors.white, size: 20),
                    ),
                  ),
                  Visibility(
                    visible: _mRecorder!.isRecording,
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
                  Visibility(
                    visible: _mplaybackReady,
                    child: GestureDetector(
                      onTap: getPlaybackFn(_mPath),
                      child: Container(
                        height: 30,
                        width: 30,
                        decoration: BoxDecoration(
                          color: Colors.black,
                          borderRadius: BorderRadius.circular(30),
                        ),
                        child: Icon(Icons.play_arrow, color: Colors.white, size: 20),
                      ),
                    ),
                  ),
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
