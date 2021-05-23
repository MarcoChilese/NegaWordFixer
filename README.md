# NegaWordFixer

This tool given in input a language ``lang``, the path where negapedia.tar.gz are stored (it considers just the latest), performs the broken words replacement in the JS variables named: ``Word2TFIDF``,
``BWord2Occur``, ``Word2Occur``.

For running the tool on docker:<br>
```
$ docker pull negapedia/negawordfixer

$ docker run -v PATH_TO_NEGAPEDIA/LANG:/out negapedia/negawordfixer --lang LANG --tar ./out
```

The available languages are:
- it
- de
- es
- fr
- en