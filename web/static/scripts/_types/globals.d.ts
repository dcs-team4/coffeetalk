// Type declarations for global variables/modules included in <script> tags.

/** Object with environment variables passed from the server through the HTML template. */
declare let serverEnv: {
  [envName: string]: string;
};

/** Diffy.js library. Minimally typed here due to lack of official type declarations. */
declare module Diffy {
  export function create(args: { onFrame: (matrix: number[][]) => void; [key: string]: any }): any;
}
