import 'package:flutter/cupertino.dart';

class Chat {
  final String ChatId;
  final int Timestamp;
  final String Title;

  Chat({
    required this.ChatId,
    required this.Timestamp,
    required this.Title,
  });
}

class ChatModel extends ChangeNotifier {
  List<Chat> _chats = [];

  bool _isNewChatCreated = false;

  bool get isNewChatCreated => _isNewChatCreated;

  List<Chat> get chats => _chats;

  void resetIsNewChatCreated() {
    _isNewChatCreated = false;
  }


  void addChat(Chat chat) {
    _chats.insert(0, chat);

    _isNewChatCreated = true;

    notifyListeners();
  }

  void setChats(List<Chat> chats) {
    _chats = chats;

    notifyListeners();
  }
}