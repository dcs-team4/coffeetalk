declare let env: {
  [envName: string]: string;
};

declare module Diffy {
  export function create(args: { onFrame: (matrix: number[][]) => void; [key: string]: any }): any;
}
