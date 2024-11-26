# DASH:
ffmpeg -i BBB-Test-Video.mp4 -map 0:v -map 0:a -f dash -seg_duration 4 \
-use_template 1 -use_timeline 1 -init_seg_name "init_$RepresentationID$.mp4" \
-media_seg_name "chunk_$RepresentationID$_$Number%05d$.m4s" \
-adaptation_sets "id=0,streams=v id=1,streams=a" ./Dash_Video/manifest.mpd
