import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:tutor/services/audio_player_service.dart';
import 'package:tutor/models/local/chat/chat_message.dart';

class MessageWidget extends StatelessWidget {
  final ChatMessage message;

  MessageWidget({Key? key, required this.message}): super(key: key);

  @override
  Widget build(BuildContext context) {
    final audioPlayerService = Provider.of<AudioPlayerService>(context);

    return Container(
      padding: EdgeInsets.only(left: 14, right: 14, top: 10, bottom: 10),
      child: Align(
        alignment: message.IsFromCurrentUser ? Alignment.topRight : Alignment.topLeft,
        child: Container(
          decoration: BoxDecoration(
            borderRadius: BorderRadius.circular(20),
            color: message.IsFromCurrentUser ? Colors.blue[200] : Colors.grey.shade200,
          ),
          padding: EdgeInsets.all(16),
          child: Column(
            children: <Widget>[
              Text(message.Text, style: TextStyle(fontSize: 15)),
              if (message.AudioUrl.isNotEmpty)
                GestureDetector(
                  onTap: () {
                    audioPlayerService.togglePlayPause(message.AudioUrl);
                  },
                  child: Container(
                    height: 30,
                    width: 30,
                    decoration: BoxDecoration(
                      color: audioPlayerService.isPlaying(message.AudioUrl) ? Colors.red : Colors.black,
                      borderRadius: BorderRadius.circular(30),
                    ),
                    child: Icon(
                      audioPlayerService.isPlaying(message.AudioUrl) ? Icons.stop : Icons.play_arrow,
                      color: Colors.white,
                      size: 20,
                    ),
                  ),
                ),
            ],
          ),
        ),
      ),
    );
  }
}
