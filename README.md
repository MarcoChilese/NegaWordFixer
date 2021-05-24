# NegaWordFixer

This tool given in input a language ``lang``, the path where negapedia.tar.gz are stored (it considers just the latest), performs the broken words replacement in the JS variables named: ``Word2TFIDF``,
``BWord2Occur``, ``Word2Occur``.

## Available Languages
The available languages are:
- it
- de
- es
- fr
- en

## Docker Run
For running the tool on docker:<br>
```
$ docker pull negapedia/negawordfixer

$ docker run -v PATH_TO_NEGAPEDIA/LANG:/out negapedia/negawordfixer --lang LANG --file "specific_file_to_process.tar.gz" --out "output_file_name.tar.gz"
```
**Flags**:
- `lang`: the language of processed Negapedia;
- `dict`: the explicit path to dictionary for building the Trie. If not specified is automatically determined through the language specified; 
  <br>
  **IMPORTANT**: the path of an external path *must* be in the same directory of the file to process, due to the directory binding of Docker.
- `out`: the explicit definition of the output filename. If not specified is set like `"fixed-"<orginal_file_name.tar.gz>`;
- `verbose`: if true, print on screen all words replacement.



