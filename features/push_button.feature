Feature: init master
  In order to push a button

  Scenario: start conversation
    When here we go again

    Given there is bot 6750404821
    When I sent "/start"
    Then there should be 3 buttons
    And the text should be:
      """
      Телеграм бот для записи онлайн

      С помощью нашего бота вы можете вести запись
      """

  Scenario:
    Given there is bot 6750404821
    When I push "🙋 Я мастер"
    Then there should be 1 buttons
      """
      Шаг 1 из 5

      Введите ваше имя

      Клиенты будут видеть ваше имя
      """
