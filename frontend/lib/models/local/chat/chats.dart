import 'package:flutter/cupertino.dart';
import 'package:tutor/models/backend/chats/get_chats_response.dart' as backend;

enum ChatType {
  General,
  JobInterview,
}

class Chat {
  final String ChatId;
  final int Timestamp;
  final String Title;
  final ChatType Type;

  Chat({
    required this.ChatId,
    required this.Timestamp,
    required this.Title,
    required this.Type,
  });

  factory Chat.fromChatResponse(backend.Chat chat) {
    return Chat(
      ChatId: chat.ChatId,
      Timestamp: chat.Timestamp,
      Title: chat.Title,
      Type: chatTypeFromInt(chat.ChatType),
    );
  }
}

ChatType chatTypeFromInt(int type) {
  switch (type) {
    case 1:
      return ChatType.General;
    case 2:
      return ChatType.JobInterview;
    default:
      return ChatType.General;
  }
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