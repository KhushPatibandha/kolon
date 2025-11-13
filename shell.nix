{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
    buildInputs = with pkgs; [
        go
        gofumpt
        goimports-reviser
        golangci-lint
        gomodifytags
        impl
        hyperfine
    ];
}

