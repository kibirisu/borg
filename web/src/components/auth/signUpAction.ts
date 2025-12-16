import { redirect, type ActionFunctionArgs } from "react-router";
import type { AppClient } from "../../lib/client";

export function signUpAction(client: AppClient) {
  return async ({ request }: ActionFunctionArgs) => {
    const formData = await request.formData();

    const username = formData.get("username")?.toString();
    const password = formData.get("password")?.toString();
    const confirmPassword = formData
      .get("confirmPassword")
      ?.toString();

    const errors: Record<string, string> = {};

    if (!username) errors.username = "Username is required";
    if (!password) errors.password = "Password is required";
    if (!confirmPassword)
      errors.confirmPassword = "Please confirm your password";

    if (password && confirmPassword && password !== confirmPassword) {
      errors.confirmPassword = "Passwords do not match";
    }

    if (Object.keys(errors).length > 0) {
      return errors;
    }

    const safeUsername = username!;
    const safePassword = password!;

    const mutation = async () => {
      return client.fetchClient.POST("/api/auth/register", {
        body: { username: safeUsername, password: safePassword },
      });
    };

    try {
      await mutation();
    } catch {
      return { form: "Registration failed" };
    }

    return redirect("/signin");
  };
}
