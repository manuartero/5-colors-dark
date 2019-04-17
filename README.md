# Simple Dark Theme

Fork from the VSCode's default dark theme but customized for simplification

## Theme

- **Text (default color)**: #E7E7E7 _almost white_
- **Background**: #151516 _almost black_

- **Keywords**: #B3DD9D _green_
- **Operators**: #8D8D8D _dark grey_
- **Functions**: #9CDCFE _light blue_
- **Classes**: #FFB1FB _pink_
- **Strings**: #D7BA7D _light orange_
- **Errors**: #F44747 _red_

### Customization

```
    "workbench.colorCustomizations": {
        "[Simple Dark Theme]": {
        }
    },
    "editor.tokenColorCustomizations": {
        "[Simple Dark Theme]": {
        }
    },
```

### Publishing

```
 $> vsce package
 $> vsce publish
```
