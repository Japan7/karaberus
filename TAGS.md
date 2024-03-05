# Tag system

This document serves as a description of the tags in karaberus.
(Ideally the tags should be described in the source code and the descriptions could
be served by karaberus itself, but it is too early for that)

# Types of tags

There are 4 types of tags:

* Generic/Kara tags: tags related to the karaoke in general.
* Media tags: tags related to the source media of the song.
* Video tags: tags related to the video of the karaoke.
* Audio tags: tags related to the audio of the karaoke.

Generic and media tags are user defined while video and music tags are
statically defined in karaberus. 

## Generic tags

Tags related to the karaoke in general.

Generic tags include:
* Titles
* Video title (needed?)
* Version (e.g. "Eurobeat version" or anything else)
* Names of the artists
* Author of the karaoke timing.
* Artists involved in the song's creation (not the video, but could be added if needed)

All these tags support alternative names/aliases.

## Media tags

Tags related to the source media.

Media tags have:
* A title
* A type of media (Anime, Game, Live action, Cartoon)
* Alternative titles if needed


## Video tags

Tags related to the video.

Video tags currently defined:
* Opening: The video is an opening video of the source media
* Ending: The video is an ending video of the source media
* Insert: The video is an insert video in the source media
* Music Video: The video is specifically made for the song
* Fanmade: The video is fanmade
* Stream: The video comes from a stream
* Concert: The video is taken from a live concert
* Advertisement: The video is an ad


## Audio tags

Tags related to the audio.

Audio tags currently defined:
* Opening: The audio is an opening song of the source media
* Ending: The audio is an ending song of the source media
* Insert: The audio is an insert song of the source media
* Live: The audio is a live interpretation of the song
