import 'package:flutter/material.dart';
import 'package:flutter_sound/flutter_sound.dart';
import 'package:logger/logger.dart';

class AudioPlayerService with ChangeNotifier {
  FlutterSoundPlayer? _player;
  String? _currentUrl;

  AudioPlayerService() {
    _player = FlutterSoundPlayer(logLevel: Level.info);
    _player!.openPlayer();
  }

  bool isPlaying(String url) {
    return _currentUrl == url && _player!.isPlaying;
  }

  void togglePlayPause(String url) async {
    if (_currentUrl == url && _player!.isPlaying) {
      await _player!.stopPlayer();
    } else {
      await _player!.startPlayer(
        fromURI: url,
        whenFinished: (){
          _player!.stopPlayer();
          _currentUrl = null;
          notifyListeners();
          },
      );
    }

    _currentUrl = url;
    notifyListeners();
  }

  @override
  void dispose() {
    _player!.closePlayer();
    super.dispose();
  }
}
