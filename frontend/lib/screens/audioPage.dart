import 'dart:io';
import 'package:audioplayers/audioplayers.dart';
import 'package:flutter/material.dart';
import 'package:record/record.dart';
import 'package:http/http.dart' as http;

class RecordingScreen extends StatefulWidget {
  const RecordingScreen({super.key});

  @override
  _RecordingScreenState createState() => _RecordingScreenState();
}

class _RecordingScreenState extends State<RecordingScreen> {
  late Record audioRecord;
  late AudioPlayer audioPlayer;
  bool isRecording = false;
  String audioPath = "";

  @override
  void initState() {
    super.initState();
    audioPlayer = AudioPlayer();
    audioRecord = Record();
  }

  @override
  void dispose() {
    super.dispose();
    audioRecord.dispose();
    audioPlayer.dispose();
  }
  bool playing=false;
  Future<void> startRecording() async {
    try {
      print("START RECODING+++++++++++++++++++++++++++++++++++++++++++++++++");
      if (await audioRecord.hasPermission()) {

        await audioRecord.start();
        setState(() {
          isRecording = true;
        });
      }
    } catch (e, stackTrace) {
      print("START RECODING+++++++++++++++++++++${e}++++++++++${stackTrace}+++++++++++++++++");
    }
  }

  Future<void> stopRecording() async {
    try {
      print("STOP RECODING+++++++++++++++++++++++++++++++++++++++++++++++++");
      String? path = await audioRecord.stop();
      setState(() {
        recoding_now=false;
        isRecording = false;
        audioPath = path!;
      });
    } catch (e) {
      print("STOP RECODING+++++++++++++++++++++${e}+++++++++++++++++++++++++++");
    }
  }

  Future<void> playRecording() async {
    try {
      playing = true;
      setState(() {});

      print("AUDIO PLAYING+++++++++++++++++++++++++++++++++++++++++++++++++");
      Source urlSource = UrlSource(audioPath);
      await audioPlayer.play(urlSource);
      // Add an event listener to be notified when the audio playback completes
      audioPlayer.onPlayerStateChanged.listen((PlayerState state) {
        if (state == PlayerState.completed) {
          playing = false;


          print("AUDIO PLAYING ENDED+++++++++++++++++++++++++++++++++++++++++++++++++");
          setState(() {});
        }
      });

    } catch (e) {
      print("AUDIO PLAYING++++++++++++++++++++++++${e}+++++++++++++++++++++++++");
    }
  }

  Future<void> pauseRecording() async {
    try {
      playing=false;

      print("AUDIO PAUSED+++++++++++++++++++++++++++++++++++++++++++++++++");

      await audioPlayer.pause();
      setState(() {

      });
      //print('Hive Playing Recording ${voiceRecordingsBox.values.cast<String>().toList().toString()}');
    } catch (e) {
      print("AUDIO PAUSED++++++++++++++++++++++++${e}+++++++++++++++++++++++++");
    }
  }

  Future<void> uploadAndDeleteRecording() async {
    try {
      final url = Uri.parse('YOUR_UPLOAD_URL'); // Replace with your server's upload URL

      final file = File(audioPath);
      if (!file.existsSync()) {
        print("UPLOADING FILE NOT EXIST+++++++++++++++++++++++++++++++++++++++++++++++++");
        return;
      }
      print("UPLOADING FILE ++++++++++++++++${audioPath}+++++++++++++++++++++++++++++++++");
      final request = http.MultipartRequest('POST', url)
        ..files.add(
          http.MultipartFile(
            'audio',
            file.readAsBytes().asStream(),
            file.lengthSync(),
            filename: 'audio.mp3', // You may need to adjust the file extension
          ),
        );

      final response = await http.Response.fromStream(await request.send());

      if (response.statusCode == 200) {
        // Upload successful, you can delete the recording if needed
        // Show a snackbar or any other UI feedback for a successful upload
        const snackBar = SnackBar(
          content: Text('Audio uploaded.'),
        );
        ScaffoldMessenger.of(context).showSnackBar(snackBar);

        // Refresh the UI
        setState(() {
          audioPath = "";
        });
      } else {
        // Handle the error or show an error message
        print('Failed to upload audio. Status code: ${response.statusCode}');
      }
    } catch (e) {
      print('Error uploading audio: $e');
    }
  }

  Future<void> deleteRecording() async {

    if (audioPath.isNotEmpty) {
      try {
        recoding_now=true;
        File file = File(audioPath);
        if (file.existsSync()) {
          file.deleteSync();
          const snackBar = SnackBar(
            content: Text('Recoding deleted'),);
          print("FILE DELETED+++++++++++++++++++++++++++++++++++++++++++++++++");
        }
      } catch (e) {
        print("FILE NOT DELETED++++++++++++++++${e}+++++++++++++++++++++++++++++++++");
      }

      setState(() {
        audioPath = "";
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        appBar: AppBar(
          title: const Text('Voice Recorder'),
        ),
        body: Row(mainAxisAlignment: MainAxisAlignment.center,
          children: [
            recoding_now? IconButton(
              icon: !isRecording
                  ? const Icon(Icons.mic_none, color: Colors.red,size: 50,)
                  : const Icon(Icons.fiber_manual_record, color: Colors.red,size: 50),
              onPressed: isRecording ? stopRecording : startRecording,
            ):
            Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                IconButton(
                  icon: !playing? Icon(Icons.play_circle, color: Colors.green,size: 50):Icon(Icons.pause_circle, color: Colors.green,size: 50),
                  onPressed: !playing?playRecording:pauseRecording,
                ),
                IconButton(
                  icon: const Icon(Icons.delete, color: Colors.red,size: 50),
                  onPressed: deleteRecording,
                ),
                IconButton(
                  icon: const Icon(Icons.trending_up, color: Colors.green,size: 50),
                  onPressed: uploadAndDeleteRecording,
                ),
              ],
            ),




          ],)






    );
  }
  bool recoding_now=true;
}