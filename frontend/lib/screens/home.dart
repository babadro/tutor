import 'package:provider/provider.dart';
import 'package:tutor/models/local/chat/chats.dart' as localChat;
import 'package:tutor/services/audio_recorder_service.dart';
import 'package:tutor/services/auth_service.dart';
import 'package:tutor/services/chat_service.dart';
import 'package:tutor/widgets/chatDetailPage.dart';
import 'package:firebase_ui_auth/firebase_ui_auth.dart';
import 'package:flutter/material.dart';
import 'package:audio_session/audio_session.dart';

class HomeScreen extends StatefulWidget {
  @override
  State<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends State<HomeScreen> {
  var selectedIndex = 0;
  localChat.Chat? selectedChat;
  late ChatService _chatService;
  AudioRecorderService? _audioRecorderService = AudioRecorderService();

  @override
  void initState() {
    super.initState();
    _chatService =
        ChatService(Provider.of<AuthService>(context, listen: false));
    _loadChats();

    _initAudioSession();
  }

  Future<void> _initAudioSession() async {
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
  }

  void _loadChats() async {
    final loadChatsRes = await _chatService.getChats();
    if (!loadChatsRes.success) {
      print('Failed to load chats: ${loadChatsRes.errorMessage}');
      // todo show error message
      return;
    }

    Provider.of<localChat.ChatModel>(context, listen: false)
        .setChats(loadChatsRes.data!);
  }

  void _deleteChat(String chatId) async {
    var deleteRes = await _chatService.deleteChat(chatId);
    if (!deleteRes.success) {
      print('Failed to delete chat: ${deleteRes.errorMessage}');
      return;
    }

    if (selectedChat?.ChatId == chatId) {
      selectedChat = null;
      selectedIndex = 0;
    }

    Provider.of<localChat.ChatModel>(context, listen: false)
        .deleteChat(chatId);
  }

  @override
  void dispose() {
    _audioRecorderService!.dispose();
    _audioRecorderService = null;
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    localChat.ChatModel chatModel = context.watch<localChat.ChatModel>();
    List<localChat.Chat> chats = chatModel.chats;
    if (chatModel.isNewChatCreated) {
      selectedIndex =
          3; // home, new generic chat, job interview chat, then old chats
      selectedChat = chats[0];
      chatModel.resetIsNewChatCreated();
    }

    List<NavigationRailDestination> getDestinations() {
      List<NavigationRailDestination> destinations = [
        NavigationRailDestination(
          icon: Icon(Icons.home),
          label: Text('Home'),
        ),
        NavigationRailDestination(
          icon: Icon(Icons.chat),
          label: Text('New Generic Chat'),
        ),
        NavigationRailDestination(
          icon: Icon(Icons.chat),
          label: Text('New Job Interview Chat'),
        ),
      ];

      // Append old chats to the destinations
      // Append old chats with three dots menu to the destinations
      /*
       destinations.addAll(chats.map((chat) => NavigationRailDestination(
            icon: Icon(Icons.chat),
            label: Text(chat.Title),
          )));
       */
      destinations.addAll(chats.map((chat) {
        return NavigationRailDestination(
          icon: Icon(Icons.chat),
          label: Row(
            mainAxisSize: MainAxisSize.min,
            children: [
              Text(chat.Title),
              PopupMenuButton<int>(
                icon: Icon(Icons.more_vert, size: 16),
                onSelected: (value) {
                  if (value == 0) {
                    _deleteChat(chat.ChatId);
                  }
                },
                itemBuilder: (context) => [
                  PopupMenuItem<int>(
                    value: 0,
                    child: Text("Delete Chat"),
                  ),
                ],
              ),
              SizedBox(width: 4), // Add some spacing between icon and text
              //Expanded(child: Text(chat.Title, overflow: TextOverflow.ellipsis)),
            ],
          ),
        );
      }));

      return destinations;
    }

    return LayoutBuilder(builder: (context, constraints) {
      return Scaffold(
        appBar: AppBar(
          actions: [
            const SignOutButton(),
            IconButton(
              icon: const Icon(Icons.person),
              onPressed: () {
                Navigator.push(
                  context,
                  MaterialPageRoute<ProfileScreen>(
                    builder: (context) => ProfileScreen(
                      appBar: AppBar(
                        title: const Text('User Profile'),
                      ),
                      actions: [
                        SignedOutAction((context) {
                          Navigator.of(context).pop();
                        })
                      ],
                      children: [
                        const Divider(),
                        Padding(
                          padding: const EdgeInsets.all(2),
                          child: AspectRatio(
                            aspectRatio: 1,
                            child: Image.asset('flutterfire_300x.png'),
                          ),
                        ),
                      ],
                    ),
                  ),
                );
              },
            )
          ],
          automaticallyImplyLeading: false,
        ),
        body: Row(
          children: <Widget>[
            SafeArea(
              child: LayoutBuilder(
                builder: (context, constraint) {
                  return SingleChildScrollView(
                    child: ConstrainedBox(
                      constraints:
                          BoxConstraints(minHeight: constraint.maxHeight),
                      child: IntrinsicHeight(
                        child: NavigationRail(
                          extended: constraints.maxWidth >= 600,
                          destinations: getDestinations(),
                          selectedIndex: selectedIndex,
                          onDestinationSelected: (value) {
                            print("Selected index: $value");
                            setState(() {
                              selectedIndex = value;

                              switch (value) {
                                case 0 || 1:
                                  selectedChat = localChat.Chat.emptyWithType(localChat.ChatType.General);
                                case 2:
                                  selectedChat = localChat.Chat.emptyWithType(localChat.ChatType.JobInterview);
                                default:
                                  selectedChat = chats[value - 3];
                              }
                            });
                          },
                        ),
                      ),
                    ),
                  );
                },
              ),
            ),
            Expanded(
              child: Container(
                color: Theme.of(context).colorScheme.primaryContainer,
                child: (selectedIndex > 0)
                    ? ChatDetailPage(
                        key: Key(selectedChat!.ChatId),
                        initialChat: selectedChat!,
                        mRecorder: _audioRecorderService!,
                      )
                    : Placeholder(),
              ),
            ),
          ],
        ),
      );
    });
  }
}
