{
  inputs = {
    systems.url = "github:nix-systems/default";
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-parts.url = "github:hercules-ci/flake-parts";
  };

  outputs = inputs @ {flake-parts, ...}:
    flake-parts.lib.mkFlake
    {
      inherit inputs;
    }
    (
      {
        inputs,
        self,
        ...
      }: {
        systems = import inputs.systems;
        perSystem = {
          pkgs,
          self',
          ...
        }: {
          packages = {
            iris = pkgs.buildGoModule {
              pname = "iris";
              version = self.shortRev or "dirty";
              src = pkgs.lib.cleanSource ./.;

              subPackages = ["cmd/iris"];

              proxyVendor = true;
              vendorHash = "sha256-kBSMhUsuCKIjAXjGfl1WSjCX+tlGi9BTnkRu9ScW6M0=";

              doCheck = false;

              meta = with pkgs.lib; {
                description = "A highly customizable, blazing fast and context-aware CLI autocomplete/navigation tool";
                homepage = "https://github.com/versenilvis/iris";
                license = licenses.mit;
                mainProgram = "iris";
              };
            };
            default = self'.packages.iris;
          };

          devShells.default = pkgs.mkShellNoCC {
            # Tell Direnv to shut up.
            DIRENV_LOG_FORMAT = "";

            packages = [
              # languages
              pkgs.go
              pkgs.pkl

              # Tooling
              pkgs.hk
              pkgs.just
              pkgs.alejandra
            ];
          };
        };
      }
    );
}
