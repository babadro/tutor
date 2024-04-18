import '../models/chat.dart';
import 'chatDetailPage.dart';
import 'package:firebase_ui_auth/firebase_ui_auth.dart';
import 'package:flutter/material.dart';

class HomeScreen extends StatefulWidget {
  @override
  State<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends State<HomeScreen> {

  var selectedIndex = 0;

  List<Chat> chats = [
    Chat(Id: '1', Title: 'Chat 1'),
    Chat(Id: '2', Title: 'Chat 2'),
    Chat(Id: '3', Title: 'Chat 3'),
  ];

  @override
  Widget build(BuildContext context) {
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


    Widget page;
    switch(selectedIndex) {
      case 0:
        page = Placeholder();
        break;
      case 1:
        page = ChatDetailPage();
        break;
      default:
        throw UnimplementedError('no widget for $selectedIndex');
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
                      });
                    },
                  ),
                ),
                Expanded(
                  child: Container(
                    color: Theme.of(context).colorScheme.primaryContainer,
                    child: page,
                  ),
                ),
              ],
            ),
          );
        }
    );
  }
}
