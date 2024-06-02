import 'package:provider/provider.dart';
import 'package:tutor/models/local/chat/chats.dart' as localChat;
import 'package:tutor/widgets/audio_page_2_flutter_sound.dart';
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
  var selectedChatId = '';
  late ChatService _chatService;

  @override
  void initState() {
    super.initState();
    _chatService = ChatService(Provider.of<AuthService>(context, listen: false));
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

    setState(() {
      Provider.of<localChat.ChatModel>(context, listen: false).
        setChats(loadChatsRes.data!);
    });
  }

  @override
  Widget build(BuildContext context) {
    List<localChat.Chat> chats = context.watch<localChat.ChatModel>().chats;
    localChat.ChatModel chatModel = context.watch<localChat.ChatModel>();
    if (chatModel.isNewChatCreated) {
      selectedIndex = 2; // 0 is home, 1 is 'new chat' button, so the first chat is at index 2
      selectedChatId = chats[0].ChatId;
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
          label: Text('New Chat'),
        ),
      ];

      // Append old chats to the destinations
      destinations.addAll(chats.map((chat) => NavigationRailDestination(
        icon: Icon(Icons.chat),
        label: Text(chat.Title),
      )));

      return destinations;
    }

    return LayoutBuilder(
        builder: (context, constraints) {
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
              children: [
                SafeArea(
                  child: NavigationRail(
                    extended: constraints.maxWidth >= 600,
                    destinations: getDestinations(),
                    selectedIndex: selectedIndex,
                    onDestinationSelected: (value) {
                      setState((){
                        selectedIndex = value;
                        selectedChatId = (value > 1) ? chats[value - 2].ChatId : '';
                      });
                    },
                  ),
                ),
                Expanded(
                  child: Container(
                    color: Theme.of(context).colorScheme.primaryContainer,
                    child: (selectedIndex > 0) ? ChatDetailPage(key: Key(selectedChatId), initialChatId: selectedChatId) : SimpleRecorder(key: Key('recording_screen')),
                  ),
                ),
              ],
            ),
          );
        }
    );
  }
}
