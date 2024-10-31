import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:tutor/services/audio_player_service.dart';
import 'package:tutor/models/local/chat/chat_message.dart';

class MessageWidget extends StatefulWidget {
  final ChatMessage message;

  MessageWidget({Key? key, required this.message}) : super(key: key);

  @override
  _MessageWidgetState createState() => _MessageWidgetState();
}

class _MessageWidgetState extends State<MessageWidget> {
  bool _isTextVisible = false;

  @override
  Widget build(BuildContext context) {
    final audioPlayerService = Provider.of<AudioPlayerService>(context);

    return Container(
      padding: EdgeInsets.only(left: 14, right: 14, top: 10, bottom: 10),
      child: Align(
        alignment: widget.message.IsFromCurrentUser ? Alignment.topRight : Alignment.topLeft,
        child: Container(
          decoration: BoxDecoration(
            borderRadius: BorderRadius.circular(20),
            color: widget.message.IsFromCurrentUser ? Colors.blue[200] : Colors.grey.shade200,
          ),
          padding: EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: <Widget>[
              // Toggle button for showing/hiding text
              TextButton(
                onPressed: () {
                  setState(() {
                    _isTextVisible = !_isTextVisible;
                  });
                },
                child: Text(
                  _isTextVisible ? "Hide Text" : "Show Text",
                  style: TextStyle(color: Colors.black54),
                ),
              ),
              // Show text only if _isTextVisible is true
              if (_isTextVisible)
                Text(
                  widget.message.Text,
                  style: TextStyle(fontSize: 15),
                ),
              // Audio play button
              if (widget.message.AudioUrl.isNotEmpty)
                IconButton(
                  onPressed: () {
                    audioPlayerService.togglePlayPause(widget.message.AudioUrl);
                  },
                  icon: Icon(
                    audioPlayerService.isPlaying(widget.message.AudioUrl) ? Icons.stop : Icons.play_arrow,
                    color: audioPlayerService.isPlaying(widget.message.AudioUrl) ? Colors.red : Colors.black,
                  ),
                  iconSize: 24, // Adjust the size as needed
                ),
            ],
          ),
        ),
      ),
    );
  }
}

