import 'package:flutter/foundation.dart' show kIsWeb;
import 'package:flutter_sound/flutter_sound.dart';
import 'package:flutter_sound_platform_interface/flutter_sound_recorder_platform_interface.dart';
import 'package:permission_handler/permission_handler.dart';
import 'package:logger/logger.dart';

const theSource = AudioSource.microphone;

class AudioRecorderService {
  Codec _codec = Codec.aacMP4;
  String _mPath = 'tau_file.mp4';
  FlutterSoundRecorder _mRecorder = FlutterSoundRecorder(logLevel: Level.info);
  bool _mRecorderIsInited = false;

  // init recorder
  void init() {
    openTheRecorder().then((_) {
      _mRecorderIsInited = true;
    });
  }

  Future<void> openTheRecorder() async {
    if (!kIsWeb) {
      var status = await Permission.microphone.request();
      if (status != PermissionStatus.granted) {
        throw RecordingPermissionException('Microphone permission not granted');
      }
    }
    await _mRecorder.openRecorder();
    if (!await _mRecorder.isEncoderSupported(_codec) && kIsWeb) {
      _codec = Codec.opusWebM;
      _mPath = 'tau_file.webm';
      if (!await _mRecorder.isEncoderSupported(_codec) && kIsWeb) {
        _mRecorderIsInited = true;
        return;
      }
    }

    _mRecorderIsInited = true;
  }

  Future<void> record() async {
    try {
      await _mRecorder.startRecorder(
        codec: _codec,
        toFile: _mPath,
        audioSource: theSource,
      );
    } catch (e) {
      Logger().e('Error recording: $e');
    }
  }

  Future<String?> stopRecording() async {
    return _mRecorder.stopRecorder();
  }

  Future<void> cancelRecording() async {
    await _mRecorder.stopRecorder();
  }

  void dispose() {
    _mRecorder.closeRecorder();
  }

  bool get inited => _mRecorderIsInited;

  bool get isStopped => _mRecorder.isStopped;

  bool get isRecording => _mRecorder.isRecording;
}
