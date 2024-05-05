import 'package:firebase_core/firebase_core.dart';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:tutor/services/auth_service.dart';

import 'app.dart';
import 'firebase_options.dart';
import 'models/local/chat/chats.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  await Firebase.initializeApp(
    options: DefaultFirebaseOptions.currentPlatform,
  );

  runApp(
      MultiProvider(
          providers: [
            Provider<AuthService>(
              create: (_) => AuthService(),
            ),
            ChangeNotifierProvider<ChatModel>(
              create: (_) => ChatModel(),
            ),
          ],
          child: MyApp()
  ));
}
