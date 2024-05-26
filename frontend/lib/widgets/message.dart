import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:tutor/services/audio_player_service.dart';

class MessageWidget extends StatelessWidget {
  final String text;
  final String audioUrl;
  final bool isFromCurrentUser;

  MessageWidget({Key? key, required this.text, required this.audioUrl, required this.isFromCurrentUser}): super(key: key);

  @override
  Widget build(BuildContext context) {
    final audioPlayerService = Provider.of<AudioPlayerService>(context);

    return Container(
      padding: EdgeInsets.only(left: 14, right: 14, top: 10, bottom: 10),
      child: Align(
        alignment: isFromCurrentUser ? Alignment.topRight : Alignment.topLeft,
        child: Container(
          decoration: BoxDecoration(
            borderRadius: BorderRadius.circular(20),
            color: isFromCurrentUser ? Colors.blue[200] : Colors.grey.shade200,
          ),
          padding: EdgeInsets.all(16),
          child: Column(
            children: <Widget>[
              Text(text, style: TextStyle(fontSize: 15)),
              if (audioUrl.isNotEmpty)
                GestureDetector(
                  onTap: () {
                    audioPlayerService.togglePlayPause(audioUrl);
                  },
                  child: Container(
                    height: 30,
                    width: 30,
                    decoration: BoxDecoration(
                      color: audioPlayerService.isPlaying(audioUrl) ? Colors.red : Colors.black,
                      borderRadius: BorderRadius.circular(30),
                    ),
                    child: Icon(
                      audioPlayerService.isPlaying(audioUrl) ? Icons.stop : Icons.play_arrow,
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
