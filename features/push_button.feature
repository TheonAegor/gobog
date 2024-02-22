Feature: init master
  In order to push a button

  Scenario: reset db
    When here we go again

  Scenario: start conversation
    Given there is bot 6750404821
    When I sent "/start"
    Then there should be 3 buttons
    And the text should be:
      """
      Телеграм бот для записи онлайн

      С помощью нашего бота вы можете вести запись
      """

  Scenario: become master
    Given there is bot 6750404821
    When I push "🙋 Я мастер"
    Then there should be 1 buttons
    And the text should be:
      """
      Шаг 1 из 5

      Введите ваше имя

      Клиенты будут видеть ваше имя
      """

  Scenario: set name
    Given there is bot 6750404821
    When I sent "Тестовое имя"
    Then there should be 7 buttons
    And the text should be:
      """
      Шаг 2 из 5

      Тестовое имя, выберите вашу специализацию?
      """

  Scenario: pick service
    Given there is bot 6750404821
    When I push "Ногтевой сервис"
    Then there should be 2 buttons
    And the text should be:
      """
      Шаг 3 из 5

      Теперь можете добавить услуги, которые вы оказываете.

      Напишите название первой услуги.
      Например, "Маникюр без покрытия"
      """

  Scenario: specify service later
    Given there is bot 6750404821
    When I push "Сделаю это позже"
    Then there should be 1 buttons
    And the text should be:
      """
      Шаг 4 из 5

      Укажите адрес, по которому вы планируете принимать клиентов.
      Если есть ориентиры, укажите их точно.

      📍 Также можете прикрепить ссылку на точку на карте.
      Пример: г. Воронеж, ул. Кольцовская, д. 15, БЦ Галерея Чижова, 2 этаж, офис 120 (справа от магазина Lime)
      """

  Scenario: set address
    Given there is bot 6750404821
    When I sent "ул. Пушкина дом Колотушкина"
    Then there should be 9 buttons
    And the text should be:
      """
      Шаг 5 из 5

      Осталось настроить ваше расписание

      Выберите ваши рабочие дни
      """

  Scenario: set schedule days
    Given there is bot 6750404821
    When I push "□ Пн"
    And I push "□ Вт"
    And I push "✅ Готово"
    Then there should be 6 buttons
    And the text should be:
      """
      Укажите ваше рабочее время.
      Далее в рамках каждого дня его можно будет изменить.
      """

  Scenario: set schedule time
    Given there is bot 6750404821
    When I push "09:00-18:00"
    Then there should be 3 buttons
    And the text should be:
      """
      Добавить перерыв?
      """

  Scenario: no breaks
    Given there is bot 6750404821
    When I push "Нет"
    Then there should be 6 buttons
    And the text should be:
      """
      На сколько дней вперед клиенты могут видеть ваше расписание?
      """

  Scenario: set schefule visibility
    Given there is bot 6750404821
    When I push "14 дней"
    Then there should be 3 buttons
    And the text should be:
      """
      Проверьте ваш график

      Запись к вам будет доступна для клиентов в: ПН, ВТ.

      Время работы в указанные дни: с 09:00 до 18:00
      Без перерыва

      Клиенты видят ваше расписание на 14 дней вперед
      """

  Scenario: confirm schedule
    Given there is bot 6750404821
    When I push "✅ Подтвердить расписание"
    Then there should be 1 image
    And there should be 1 buttons
    And the text should be:
      """
      Тестовое имя, ваш профиль готов!
      К вам открыта запись по ПН, ВТ

      Время работы:
      c 09:00 до 18:00
      без перерыва

      Теперь вы можете разместить ссылку для записи. Как только клиент запишется, мы пришлем уведомление и добавим запись в ваш календарь
      https://t.me/test_beauty36_bot?start=1
      """

  Scenario: go back to main page
    Given there is bot 6750404821
    When I push "🏡 На главную"
    Then there should be 5 buttons
    And the text should be:
      """
      Тестовое имя, к вам пока никто не записался 🙃

      Вы можете разместить ссылку для записи. Как только клиент запишется, мы пришлем уведомление и добавим запись в ваш календарь
      https://t.me/test_beauty36_bot?start=1
      """

  Scenario: to MAIN page
    Given there is bot 6750404821
    When I push "👩‍🦳 Beauty 36"
    Then there should be 3 buttons
    And the text should be:
      """
      Телеграм бот для записи онлайн

      С помощью нашего бота вы можете вести запись
      """