import 'dart:convert';

import 'package:provider/provider.dart';
import '../models/backend/chats/get_chats_response.dart';
import '../services/auth_service.dart';
import 'chatDetailPage.dart';
import 'package:firebase_ui_auth/firebase_ui_auth.dart';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;

class HomeScreen extends StatefulWidget {
  @override
  State<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends State<HomeScreen> {
  var selectedIndex = 0;
  List<Chat> chats = [];
  String selectedChatId = '';

  @override
  void initState() {
    super.initState();
    _loadChats();
  }

  void _loadChats() async {
    const apiUrl = 'http://localhost:8080/chats?limit=100&timestamp=0';
    final uri = Uri.parse(apiUrl);

    final authService = Provider.of<AuthService>(context, listen: false);

    String? authToken = await authService.getCurrentUserIdToken();

    try {
      print('Fetching chats from $uri');
      final response = await http.get(
        uri,
        headers: {
          'Authorization': 'Bearer $authToken', // Include the authorization header
          'Content-Type': 'application/json',
        },
      ).timeout(Duration(seconds: 10));
      if (response.statusCode == 200) {
        final chatsResponse = GetChatsResponse.fromJson(jsonDecode(response.body) as Map<String, dynamic>);

        setState(() {
          this.chats = chatsResponse.Chats;
        });
      } else {
        print('Failed to fetch chats: ${response.statusCode}');
      }
    } catch (e) {
      print('Failed to fetch chats: $e');
    }
  }

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
        label: Text(chat.ChatId), // todo: get the name of the chat
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
                    child: (selectedIndex > 0) ? ChatDetailPage(key: Key(selectedChatId), chatId: selectedChatId) : Placeholder(),
                  ),
                ),
              ],
            ),
          );
        }
    );
  }
}
