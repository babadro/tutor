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
  int CurrentMessageIDx;

  Chat({
    required this.ChatId,
    required this.Timestamp,
    required this.Title,
    required this.Type,
    this.CurrentMessageIDx = 0,
  });

  factory Chat.fromChatResponse(backend.Chat chat) {
    return Chat(
      ChatId: chat.ChatId,
      Timestamp: chat.Timestamp,
      Title: chat.Title,
      Type: chatTypeFromInt(chat.ChatType),
      CurrentMessageIDx: chat.CurrentMessageIDx,
    );
  }

  Chat.emptyWithType(ChatType type)
      : ChatId = '',
        Timestamp = 0,
        Title = '',
        Type = type,
        CurrentMessageIDx = 0;
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

int chatTypeToInt(ChatType type) {
  switch (type) {
    case ChatType.JobInterview:
      return 2;
    default:
      return 1;
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

  void deleteChat(String chatId) {
    _chats.removeWhere((chat) => chat.ChatId == chatId);

    notifyListeners();
  }
}
