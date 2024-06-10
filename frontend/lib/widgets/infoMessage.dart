import 'package:flutter/material.dart';

class InfoMessageWidget extends StatelessWidget {
  final bool isRecording;

  InfoMessageWidget({Key? key, required this.isRecording}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Align(
      alignment: Alignment.topRight,
      child: Padding(
        padding: const EdgeInsets.all(20),
        child: Container(
          decoration: BoxDecoration(
            borderRadius: BorderRadius.circular(20),
            color: isRecording ? Colors.orange : Colors.blue,
          ),
          padding: EdgeInsets.all(16),
          child: Row(
            mainAxisSize: MainAxisSize.min,
            children: [
              Icon(isRecording ? Icons.mic : Icons.send, color: Colors.white),
              SizedBox(width: 8),
              Text(
                isRecording ? 'Recording' : 'Sending',
                style: TextStyle(color: Colors.white),
              ),
              SizedBox(width: 8),
              CircularProgressIndicator(),
            ],
          ),
        ),
      ),
    );
  }
}