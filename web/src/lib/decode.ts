export interface DecodedToken {
  username: string | null;
  userId: string | null;
  issuer?: string | null;
}

const decodeToken = (token: string | null): DecodedToken | null => {
  if (!token) {
    return null;
  }
  const raw = token.replace(/^Bearer:\s*/i, "");
  const base64Url = raw.split(".")[1];
  if (!base64Url) {
    return null;
  }
  const base64 = base64Url.replace(/-/g, "+").replace(/_/g, "/");
  const jsonPayload = decodeURIComponent(
    atob(base64)
      .split("")
      .map((c) => `%${`00${c.charCodeAt(0).toString(16)}`.slice(-2)}`)
      .join(""),
  );
  const payload = JSON.parse(jsonPayload) as {
    username?: string;
    name?: string;
    sub?: number | string;
    iss?: string;
  };
  const sub = payload.sub;
  const userId =
    sub === undefined || sub === null
      ? null
      : typeof sub === "string"
        ? sub
        : String(sub);
  return {
    username: payload.username ?? payload.name ?? null,
    userId,
    issuer: payload.iss ?? null,
  };
};

export default decodeToken;
