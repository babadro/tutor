import 'package:flutter/material.dart';
import 'package:just_audio/just_audio.dart';

class AudioPlayerService with ChangeNotifier {
  AudioPlayer _player = AudioPlayer();
  String? _currentUrl;

  bool isPlaying(String url) {
    return _currentUrl == url;
  }

  void togglePlayPause(String url) async {
    if (_currentUrl == url) {
      await _player.stop();
      _currentUrl = null;
      notifyListeners();

      return;
    }

    await _player.stop();
    await _player.setUrl(url);

    _currentUrl = url;
    notifyListeners();

    await _player.play();
    _currentUrl = null;


    notifyListeners();
  }

  @override
  void dispose() {
    _player.dispose();
    super.dispose();
  }
}
