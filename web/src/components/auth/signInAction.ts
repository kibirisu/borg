import type { ActionFunctionArgs } from "react-router";
import { redirect } from "react-router";
import type { AppClient } from "../../lib/client";

export function signInAction(client: AppClient) {
  return async ({ request }: ActionFunctionArgs) => {
    const formData = await request.formData();

    const username = formData.get("username")?.toString();
    const password = formData.get("password")?.toString();

    const errors: Record<string, string> = {};

    if (!username) errors.username = "Username is required";
    if (!password) errors.password = "Password is required";

    if (Object.keys(errors).length > 0) {
      return errors;
    }

    const safeUsername = username!;
    const safePassword = password!;

    const mutation = async () => {
      return client.fetchClient.POST("/api/auth/login", {
        body: { username: safeUsername, password: safePassword },
      });
    };

    try {
      await mutation();
    } catch {
      return { form: "Invalid username or password" };
    }

    return redirect("/");
  };
}
