{
  description = "IRIS - Shell commands suggestion tool";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs = { self, nixpkgs }:
    let
      supportedSystems = [ "x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin" ];
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;
      pkgsFor = system: nixpkgs.legacyPackages.${system};
    in
    {
      packages = forAllSystems (system:
        let
          pkgs = pkgsFor system;
        in
        {
          iris = pkgs.buildGoModule {
            pname = "iris";
            version = self.shortRev or "dirty";
            src = ./.;

            subPackages = [ "cmd/iris" ];

            proxyVendor = true;
            vendorHash = "sha256-kBSMhUsuCKIjAXjGfl1WSjCX+tlGi9BTnkRu9ScW6M0=";

            doCheck = false;

            meta = with pkgs.lib; {
              description = "A shell auto-completion tool that works like code editor's IntelliSense";
              homepage = "https://github.com/versenilvis/iris";
              license = licenses.mit;
              mainProgram = "iris";
            };
          };
          default = self.packages.${system}.iris;
        });

      devShells = forAllSystems (system:
        let
          pkgs = pkgsFor system;
        in
        {
          default = pkgs.mkShell {
            buildInputs = with pkgs; [ go ];
          };
        });
    };
}
